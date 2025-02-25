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
const PERM_CODE = 0644

func createTask(desc string) {
	file, err := os.OpenFile(DB_FILE, os.O_RDWR|os.O_CREATE, PERM_CODE)
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

	err = os.WriteFile(DB_FILE, data, PERM_CODE)

	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("Output: Task added successfully (ID: %d)\n", idCounter)
}

func updateTask(id int, desc string) {
	file, err := os.OpenFile(DB_FILE, os.O_RDWR, PERM_CODE)

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

	err = os.WriteFile(DB_FILE, data, PERM_CODE)

	if err != nil {
		panic(err.Error())
	}
}

func deleteTask(id int) {
	file, err := os.OpenFile(DB_FILE, os.O_RDWR, PERM_CODE)

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
			list = append(list[:i], list[i+1:]...)
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

	err = os.WriteFile(DB_FILE, data, PERM_CODE)

	if err != nil {
		panic(err.Error())
	}
}

func changeStatus(id int, status string) {
	file, err := os.OpenFile(DB_FILE, os.O_RDWR, PERM_CODE)

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
			list[i].Status = status
			list[i].UpdateAt = time.Now()
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

	err = os.WriteFile(DB_FILE, data, PERM_CODE)

	if err != nil {
		panic(err.Error())
	}
}

func list(filter string) {
	data, err := os.ReadFile(DB_FILE)

	if err != nil {
		panic(err.Error())
	}

	list := []Task{}
	err = json.Unmarshal(data, &list)

	if err != nil {
		panic(err.Error())
	}

	if filter == "done" || filter == "todo" || filter == "in-progress" {
		for _, task := range list {
			if task.Status == filter {
				fmt.Printf("%d %s\n", task.Id, task.Description)
			}
		}

		return
	}

	for _, task := range list {
		fmt.Printf("%d %s\n", task.Id, task.Description)
	}
}

func readCommand(action *string, taskId *int, taskText *string, status *string, listFilter *string) error {
	for i, arg := range os.Args {

		if i == 1 {
			*action = arg
		} else if *action == "add" {

			if i == 2 {
				*taskText = arg
			}

		} else if *action == "delete" {

			if i == 2 {
				argInt, err := strconv.ParseInt(arg, 10, 64)

				if err != nil {
					return fmt.Errorf("invalid task id: %s", arg)
				}

				*taskId = int(argInt)
			} else if i == 3 {
				*taskText = arg
			}

		} else if *action == "update" {

			if i == 2 {
				argInt, err := strconv.ParseInt(arg, 10, 64)

				if err != nil {
					return fmt.Errorf("invalid task id: %s", arg)
				}

				*taskId = int(argInt)
			} else if i == 3 {
				*taskText = arg
			}

		} else if *action == "mark-in-progress" || *action == "mark-done"  {

			if *action == "mark-in-progress" {
				*status = "in-progress"
			} else {
				*status = "done"
			}

			if i == 2 {
				argInt, err := strconv.ParseInt(arg, 10, 64)

				if err != nil {
					return fmt.Errorf("invalid task id: %s", arg)
				}

				*taskId = int(argInt)
			}

		} else if *action == "list" {

			if i == 2 {
				*listFilter = arg
			}

		}
	}

	return nil
}

func main() {
	var (
		action   string
		taskId   int
		taskText string
		status string
		listFilter string
	)

	err := readCommand(&action, &taskId, &taskText, &status, &listFilter)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if action == "add" {
		createTask(taskText)
	} else if action == "update" {
		updateTask(taskId, taskText)
	} else if action == "delete" {
		deleteTask(taskId)
	} else if action == "mark-in-progress" || action == "mark-done" {
		changeStatus(taskId, status)
	} else if action == "list" {
		list(listFilter)
	}
}
