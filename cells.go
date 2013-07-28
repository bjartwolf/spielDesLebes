package main
import "time"

type Cell struct {
    x int
    y int
    levande bool
    neighbors chan string
    subscribers []chan string
 }

// Subscribe adds a channel to a cells list of subscribers for notification 
func (c *Cell) Subscribe(subscriber chan string) { // could return dispose method to unsubscribe like rx?  
    c.subscribers = append(c.subscribers,subscriber)
}

// Notify notifies subscribing cells with a message with a delay [ms]
func (c *Cell) notify(msg string, delay time.Duration) {
    go func(c *Cell, msg string) {
            time.Sleep(delay*time.Millisecond)
            for _, s := range c.subscribers {
                s <- msg
            }
    }(c, msg)
}

// Notify notifies subscribing cells with a message with a delay [ms]
func (c *Cell) døy() {
    c.levande = false
    c.notify("fe đøyr, frendar døyr, en sjølv døyr på samme vis", 500)
}

func (c *Cell) våkne() {
    c.levande = true
    c.notify("vi ble født i samme natt, og det blir ikke langt mellom vår død heller", 500)
}

func (c *Cell) Spawn() {
         levandeGranner := 0
         for {
             select {
                 case msg := <-c.neighbors:
                    switch msg {
                        case "fe đøyr, frendar døyr, en sjølv døyr på samme vis":
                            levandeGranner--
                        case "vi ble født i samme natt, og det blir ikke langt mellom vår død heller":
                            levandeGranner++
                    }
                 case <- time.Tick(time.Second):
                    if (!c.levande && levandeGranner== 3) {
                        c.våkne()
                    } else if (c.levande && (levandeGranner < 2 || levandeGranner > 3)) {
                        c.døy()
                    }
               }
          }
}

