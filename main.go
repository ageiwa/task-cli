package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

type Task struct {
	Id          int       `json:"id"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdateAt    time.Time `json:"updateAt"`
}

const DB_FILE = "db.json"

func createTask(desc string) {
	file, err := os.OpenFile(DB_FILE, os.O_RDWR|os.O_CREATE, 0644)
	idCounter := 0

	if err != nil {
		panic(err.Error())
	}

	defer file.Close()

	list := []Task{}
	err = json.NewDecoder(file).Decode(&list)

	if err != nil && err != io.EOF {
		panic(err.Error())
	}

	if len(list) > 0 {
		lastIndex := len(list) - 1
		idCounter = list[lastIndex].Id + 1
	} else {
		idCounter += 1
	}

	task := Task{
		Id:          idCounter,
		Description: desc,
		Status:      "todo",
		CreatedAt:   time.Now(),
	}

	list = append(list, task)
	data, err := json.Marshal(list)

	if err != nil && err == io.EOF {
		fmt.Println("Файл пуст")
		return
	}

	_, err = file.WriteAt(data, 0)

	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("Output: Task added successfully (ID: %d)\n", idCounter)
}

func updateTask(id int, desc string) {
	file, err := os.OpenFile(DB_FILE, os.O_RDWR, 0644)

	if err != nil {
		panic(err.Error())
	}

	defer file.Close()

	list := []Task{}
	err = json.NewDecoder(file).Decode(&list)

	if err != nil {
		panic(err.Error())
	}

	isFind := false

	for i, task := range list {
		if task.Id == id {
			list[i].Description = desc
			list[i].UpdateAt = time.Now()
			isFind = true
			break
		}
	}

	if !isFind {
		fmt.Println("No task with id:", id)
		return
	}

	data, err := json.Marshal(list)

	if err != nil {
		panic(err.Error())
	}

	_, err = file.WriteAt(data, 0)

	if err != nil {
		panic(err.Error())
	}
}

func markInProgress(id int) {
	file, err := os.OpenFile(DB_FILE, os.O_RDWR, 0644)

	if err != nil {
		panic(err.Error())
	}

	list := []Task{}
	err = json.NewDecoder(file).Decode(&list)

	if err != nil {
		panic(err.Error())
	}

	isFind := false

	for i, task := range list {
		if task.Id == id {
			list[i].Status = "in progress"
			isFind = true
			break
		}
	}

	if !isFind {
		panic("No task with id: ")
	}

	data, err := json.Marshal(list)

	if err != nil {
		panic(err.Error())
	}

	_, err = file.WriteAt(data, 0)

	if err != nil {
		panic(err.Error())
	}
}

func list() {
	file, err := os.ReadFile(DB_FILE)

	if err != nil {
		panic(err.Error())
	}

	list := []Task{}
	err = json.Unmarshal(file, &list)

	if err != nil {
		panic(err.Error())
	}

	for _, task := range list {
		fmt.Printf("%d %s\n", task.Id, task.Description)
	}
}

func readCommand(action *string, taskId *int, taskText *string) error {
	for i, arg := range os.Args {

		if i == 1 {
			*action = arg
		} else if *action == "add" {

			if i == 2 {
				*taskText = arg
			}

		} else if *action == "update" {

			if i == 2 {
				argInt, err := strconv.ParseInt(arg, 10, 64)

				if err != nil {
					return fmt.Errorf("Invalid task id: %s", arg)
				}

				*taskId = int(argInt)
			} else if i == 3 {
				*taskText = arg
			}

		} else if *action == "mark-in-progress" {

			if i == 2 {
				argInt, err := strconv.ParseInt(arg, 10, 64)

				if err != nil {
					return fmt.Errorf("Invalid task id: %s", arg)
				}

				*taskId = int(argInt)
			}

		}
	}

	return nil
}

func main() {
	list()

	// var (
	// 	action   string
	// 	taskId   int
	// 	taskText string
	// )

	// err := readCommand(&action, &taskId, &taskText)

	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }

	// if action == "add" {
	// 	createTask(taskText)
	// } else if action == "update" {
	// 	updateTask(taskId, taskText)
	// } else if action == "mark-in-progress" {
	// 	markInProgress(taskId)
	// }
}
