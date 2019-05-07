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
        "net/http"
        "os"
        "io"
        "bytes"
        "strings"
        "errors"
)

//Function DownloadFile() copy from :
//https://golangcode.com/download-a-file-from-a-url/
func DownloadFile(filepath string, url string) error {
    // Create the file
    out, err := os.Create(filepath)
    if err != nil {
        return err
    }
    defer out.Close()
    // Get the data
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    // Write the body to file
    _, err = io.Copy(out, resp.Body)
    if err != nil {
        return err
    }
    return nil
}

func UpdateDB(file string) map[string][]map[string]interface{} {
  err := DownloadFile("rollup_db.yaml","https://raw.githubusercontent.com/Sellto/K8S_AutoLabeler/master/db.yaml")
  if err != nil {
      log.Printf("update db error:",err)
  }
  db := make(map[string][]map[string]interface{})
  var yamlFile []byte
  yamlFile, err = ioutil.ReadFile("rollup_db.yaml")
  if err != nil {
    log.Printf("downloaded file cannot be read")
  }
  err = yaml.Unmarshal([]byte(yamlFile), &db)
  if err != nil {
    log.Printf("error when try to parse the downloaded database")
    log.Printf("try to read a older local version")
    yamlFile, err = ioutil.ReadFile(file)
    if err != nil {
        log.Printf("can't find the db.yaml file")
    }
    err = yaml.Unmarshal([]byte(yamlFile), &db)
    if err != nil {
      log.Fatalf("check the db.yaml format")
    }
  } else {
    os.Remove("db.yaml")
    os.Rename("rollup_db.yaml",file)
  }
  //Parse into a map of list of maps
  return db
}


func GetLabel(device *udev.Device) (bool,[]string) {
  db := UpdateDB("db.yaml")
  log.Println("device :",device.PropertyValue("ID_VENDOR"),"-",device.PropertyValue("ID_MODEL"))
  log.Println("path:",device.PropertyValue("DEVNAME"))
  for key, value := range db {
    if key == device.PropertyValue("ID_VENDOR_ID")+":"+device.PropertyValue("ID_MODEL_ID"){
      for _,description := range value {
        if description["Manufacturer"] == device.PropertyValue("ID_VENDOR") && description["Product"] == device.PropertyValue("ID_MODEL") {
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
}


func Pipe(c1,c2 *exec.Cmd) string {
  var out bytes.Buffer
  r1, w1 := io.Pipe()
  c1.Stdout = w1
  c2.Stdin = r1
  c2.Stdout = &out
  c1.Start()
  c2.Start()
  c1.Wait()
  w1.Close()
  c2.Wait()
  return out.String()
}

func extractNodeHost(stdout string) (string,error) {
  splited_stdout := strings.Split(stdout," ")
  cleaned_stdout := []string{}
  for _,value := range splited_stdout {
    if value != "" {
      cleaned_stdout = append(cleaned_stdout,value)
    }
  }
  if len(cleaned_stdout) > 6 {
    return cleaned_stdout[6],nil
  } else {
    return "no way",errors.New("Can't Extract data")
  }
}

func LabeliseNode(device *udev.Device, label []string) {
  hostname, _ := os.Hostname()
  c1 := exec.Command("kubectl","get","pod","-o","wide")
  c2 := exec.Command("grep", hostname)
  result := Pipe(c1,c2)
  nodehost,err := extractNodeHost(result)
  if err != nil {
      log.Println(err)
  }
  endofcmd := "=yes"
  for _ , lab := range label {
    if device.Action() == "remove" {
      endofcmd = "-"
    }
    cmd := exec.Command("kubectl", "label","nodes",nodehost,lab+endofcmd)
    out, err := cmd.CombinedOutput()
    if err != nil {
      log.Print("label command failed with",string(out))
    } else {
      log.Println(device.Action(),"label:",lab)
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
