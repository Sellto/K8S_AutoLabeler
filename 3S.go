package main

import (
        "fmt"
        "log"
        "io/ioutil"
        "gopkg.in/yaml.v2"
        "github.com/jochenvg/go-udev"
        "context"
        "sync"
        "os/exec"
)

func UpdateDB(file string) map[string][]map[string]interface{} {
  //Read yaml file.
  yamlFile, err := ioutil.ReadFile(file)
  if err != nil {
    log.Printf("yamlFile.Get err   #%v ", err)
  }
  //Parse into a map of list of maps
  db := make(map[string][]map[string]interface{})
  err = yaml.Unmarshal([]byte(yamlFile), &db)
  if err != nil {
    log.Fatalf("error: %v", err)
  }
  return db
}


func GetLabel(device *udev.Device) (bool,[]string) {
  db := UpdateDB("db.yaml")
  log.Println("device :",device.PropertyValue("ID_VENDOR"),"-",device.PropertyValue("ID_MODEL"))
  log.Println("path:",d.PropertyValue("DEVNAME"))
  for key, value := range db {
    if key == d.PropertyValue("ID_VENDOR_ID")+":"+device.PropertyValue("ID_MODEL_ID"){
      for _,description := range value {
        if description["Manufacturer"] == device.PropertyValue("ID_VENDOR") && description["Product"] == d.PropertyValue("ID_MODEL") {
            label_db := description["Label"].([]interface{})
            label_list := make([]string, len(label_db))
            for i, v := range label_db {
                label_list[i] = fmt.Sprint(v)
            }
            return true, label_list
        }
      }
    }
  }
  return false, []string{}
}d

func LabeliseNode(device *udev.Device, label []string) {
  endofcmd := "=yes"
  for _ , lab := range label {
    if device.Action() == "remove" {
      endofcmd = "-"
    }
    cmd := exec.Command("kubectl", "label","nodes","masterone",lab+endofcmd)
    out, err := cmd.CombinedOutput()
    if err != nil {
      log.Print("label command failed with",string(out))
    } else {
      log.Println(d.Action(),"label:",lab)
    }
  }
}


func main() {
  myudev := udev.Udev{}
  enum := myudev.NewEnumerate()
  fmt.Println("Already plugged interface :")
  enum.AddMatchSubsystem("tty")
  enum.AddMatchIsInitialized()
  devices, _ := enum.Devices()
  for _,device := range devices {
    parent := device.ParentWithSubsystemDevtype("usb", "usb_device")
    if parent != nil {
      {
        if dev_in_db, label_list := GetLabel(device); dev_in_db {
          LabeliseNode(device,label_list)
        } else  {
          log.Println("device information not found in the database")
        }
      }
    }
  }


  //Monitor the plug and deplug radio device
  m := myudev.NewMonitorFromNetlink("udev")
  // Add filters to monitor
  m.FilterAddMatchSubsystem("tty")
	// Create a context
	ctx, _ := context.WithCancel(context.Background())
	// Start monitor goroutine and get receive channel
	channel, _ := m.DeviceChan(ctx)
	// WaitGroup for timersFatalf
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		fmt.Println("\nStarted listening on channel :")
		for device := range channel {
      if dev_in_db, label_list := GetLabel(device); dev_in_db {
          LabeliseNode(device,label_list)
        } else  {
          log.Println("device information not found in the database")
        }
		}
		fmt.Println("Channel closed")
		wg.Done()
	}()
	wg.Wait()
}
