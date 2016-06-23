package generator

import (
  "io"
  "fmt"
  "text/template"

  "github.com/docker/docker/builder/dockerfile/parser"
	"bytes"
)


var nomadTemplate = `
job "{{.JobName}}" {
	datacenters = ["dc1"]
	type = "service"

	# Configure the job to do rolling updates
	update {
		# Stagger updates every 10 seconds
		stagger = "10s"

		# Update a single task at a time
		max_parallel = 1
	}

	group "default" {
		count = 1

		restart {
			attempts = 10
			interval = "5m"
			delay = "25s"
			mode = "delay"
		}

		task "{{.TaskName}}" {
			# Use Docker to run the task.
			driver = "docker"

			# Configure Docker driver with the image
			config {
				image = "{{.Image}}"
			}
      
      {{range $element := .Ports}}
			service {
				name = "${TASKGROUP}-{{$element}}"
				port = {{$element}}
				check {
					name = "alive"
					type = "tcp"
					interval = "10s"
					timeout = "2s"
				}
			}

      {{end}}

			# We must specify the resources required for
			# this task to ensure it runs on a machine with
			# enough capacity.
			resources {
				cpu = 500 # 500 Mhz
				memory = 256 # 256MB
			}
		}
	}
}
`


type nomadInfo struct {
  JobName string;
  Image string
  Ports []string
  TaskName string
}

func GenerateFromReader(input io.Reader, jobName string) ([]byte, error) {
  node, err := parser.Parse(input)
  if err != nil {
    return []byte{}, nil
  }

  info := &nomadInfo{
    Ports: make([]string, 0),
    JobName: jobName,
  }

  fmt.Println(node.Next)

  for _, child := range node.Children {
    switch {
      case child.Value == "expose" && child.Next != nil:
        info.Ports = append(info.Ports, child.Next.Value)
      break;
      case child.Value == "from" && child.Next != nil:
        info.Image = child.Next.Value
        info.TaskName = child.Next.Value
      break;
    }
  }

  b := bytes.NewBuffer([]byte{})
  t := template.Must(template.New("nomad").Parse(nomadTemplate))
  err = t.Execute(b, info)
  if err != nil {
    return []byte{}, err
  }

  return b.Bytes(), nil
}