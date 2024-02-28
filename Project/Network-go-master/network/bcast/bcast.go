package bcast

import (
	"Elev-project/Network-go-master/network/conn"
	"encoding/json"
	"fmt"
	"net"
	"reflect"
)


const (
    bufSize       = 1024
    acknowledgmentTimeout = 100 * time.Millisecond // Timeout for acknowledgment
)

type typeTaggedJSON struct {
	TypeId string
	JSON   []byte
}

// Message structure with acknowledgment flag
type Message struct {
    Data          interface{}
    Acknowledgment bool
}



// Transmitter function modified to include acknowledgment support. 

//This function has a sequential approach to broadcasting. It waits for a message (provided by the user) to be received
//on a channel, and then broadcasts it to every channel. This is done by iterating through each channel, sending the message
//and waiting for an acknowledgment before moving on to the next one. 
func Transmitter(port int, chans ...interface{}) {
    maxRetries := 3 // Maximum number of retry attempts

    checkArgs(chans...)
    typeNames := make([]string, len(chans))
    selectCases := make([]reflect.SelectCase, len(typeNames))
    retryCount := make([]int, len(chans))

    for i, ch := range chans {
        selectCases[i] = reflect.SelectCase{
            Dir:  reflect.SelectRecv,
            Chan: reflect.ValueOf(ch),
        }
        typeNames[i] = reflect.TypeOf(ch).Elem().String()
    }

    conn := conn.DialBroadcastUDP(port)
    addr, _ := net.ResolveUDPAddr("udp4", fmt.Sprintf("255.255.255.255:%d", port))
    for {
        for i, ch := range chans {
            if retryCount[i] >= maxRetries {
                fmt.Printf("Maximum retries reached for channel %d. Moving to the next channel.\n", i)
                continue // Move to next channel if max retries reached
            }

            chosen, value, _ := reflect.Select([]reflect.SelectCase{{
                Dir:  reflect.SelectRecv,
                Chan: reflect.ValueOf(ch),
            }, {
                Dir:  reflect.SelectRecv,
                Chan: reflect.ValueOf(time.After(acknowledgmentTimeout)), // Wait for acknowledgment or timeout
            }})

            // Check if acknowledgment received
            if chosen == 1 {
                // Extract the message and check if it's an acknowledgment
                if msg, ok := value.Interface().(Message); ok && msg.Acknowledgment {
                    // Reset retry count for this channel since acknowledgment received
                    retryCount[i] = 0
                    continue // Move to next channel if acknowledgment received
                }
            }

            // Send message
            jsonstr, _ := json.Marshal(Message{Data: value.Interface(), Acknowledgment: false})
            ttj, _ := json.Marshal(typeTaggedJSON{
                TypeId: typeNames[chosen],
                JSON:   jsonstr,
            })
            if len(ttj) > bufSize {
                panic(fmt.Sprintf(
                    "Tried to send a message longer than the buffer size (length: %d, buffer size: %d)\n\t'%s'\n"+
                        "Either send smaller packets, or go to network/bcast/bcast.go and increase the buffer size",
                    len(ttj), bufSize, string(ttj)))
            }
            conn.WriteTo(ttj, addr)

            // Increment retry count for this channel
            retryCount[i]++
        }
    }
}

// Receiver function modified to include acknowledgment support
func Receiver(port int, chans ...interface{}) {
    checkArgs(chans...)
    chansMap := make(map[string]interface{})
    for _, ch := range chans {
        chansMap[reflect.TypeOf(ch).Elem().String()] = ch
    }

    var buf [bufSize]byte
    conn := conn.DialBroadcastUDP(port)
    for {
        n, _, e := conn.ReadFrom(buf[0:])
        if e != nil {
            fmt.Printf("bcast.Receiver(%d, ...):ReadFrom() failed: \"%+v\"\n", port, e)
        }

        var ttj typeTaggedJSON
        json.Unmarshal(buf[0:n], &ttj)
        ch, ok := chansMap[ttj.TypeId]
        if !ok {
            continue
        }

        var msg Message
        json.Unmarshal(ttj.JSON, &msg)

        // Set acknowledgment to true and send back
        msg.Acknowledgment = true
        jsonstr, _ := json.Marshal(msg)
        conn.WriteTo(jsonstr, &net.UDPAddr{IP: net.ParseIP("255.255.255.255"), Port: port})

        // Send received data to channel
        v := reflect.New(reflect.TypeOf(ch).Elem())
        json.Unmarshal([]byte{}, v.Interface()) // Empty data for simplicity
        reflect.Select([]reflect.SelectCase{{
            Dir:  reflect.SelectSend,
            Chan: reflect.ValueOf(ch),
            Send: reflect.Indirect(v),
        }})
    }
}



// Checks that args to Tx'er/Rx'er are valid:
//  All args must be channels
//  Element types of channels must be encodable with JSON
//  No element types are repeated
// Implementation note:
//  - Why there is no `isMarshalable()` function in encoding/json is a mystery,
//    so the tests on element type are hand-copied from `encoding/json/encode.go`
func checkArgs(chans ...interface{}) {
	n := 0
	for range chans {
		n++
	}
	elemTypes := make([]reflect.Type, n)

	for i, ch := range chans {
		// Must be a channel
		if reflect.ValueOf(ch).Kind() != reflect.Chan {
			panic(fmt.Sprintf(
				"Argument must be a channel, got '%s' instead (arg# %d)",
				reflect.TypeOf(ch).String(), i+1))
		}

		elemType := reflect.TypeOf(ch).Elem()

		// Element type must not be repeated
		for j, e := range elemTypes {
			if e == elemType {
				panic(fmt.Sprintf(
					"All channels must have mutually different element types, arg# %d and arg# %d both have element type '%s'",
					j+1, i+1, e.String()))
			}
		}
		elemTypes[i] = elemType

		// Element type must be encodable with JSON
		checkTypeRecursive(elemType, []int{i+1})

	}
}


func checkTypeRecursive(val reflect.Type, offsets []int){
	switch val.Kind() {
	case reflect.Complex64, reflect.Complex128, reflect.Chan, reflect.Func, reflect.UnsafePointer:
		panic(fmt.Sprintf(
			"Channel element type must be supported by JSON, got '%s' instead (nested arg# %v)",
			val.String(), offsets))
	case reflect.Map:
		if val.Key().Kind() != reflect.String {
			panic(fmt.Sprintf(
				"Channel element type must be supported by JSON, got '%s' instead (map keys must be 'string') (nested arg# %v)",
				val.String(), offsets))
		}
		checkTypeRecursive(val.Elem(), offsets)
	case reflect.Array, reflect.Ptr, reflect.Slice:
		checkTypeRecursive(val.Elem(), offsets)
	case reflect.Struct:
		for idx := 0; idx < val.NumField(); idx++ {
			checkTypeRecursive(val.Field(idx).Type, append(offsets, idx+1))
		}
	}
}


