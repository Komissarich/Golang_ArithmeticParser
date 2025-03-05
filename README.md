# Асинхронный арифметический калькулятор для курса по Golang от Яндекса

## Содержание

- [Установка](#установка)
- [Запуск](#запуск)
- [Работа](#Работа)
- [Конфигурирование](#Конфигурирование)
- [Тестирование](#Тестирование)


## Установка

Склонируйте репозиторий на локальный компьютер:

```bash
git clone https://github.com/Komissarich/Golang_ArithmeticParser
cd Golang_ArithmeticParser
```

## Запуск

```bash
go run ./cmd
```
**Зайдите в ваш браузер и введите localhost:8080**

## Работа
Вы можете вводить арифметические выражения используя кнопки на калькуляторе. После набора выражения, нажмите =, для отправки выражения на сервер.
Вы можете убедиться, что сервер получил ваш запрос, если нажмете Receive с указанием All expression. 

![image](https://github.com/user-attachments/assets/77b1c8a7-a650-4c54-a21d-d7ba1ea2ef86)

Если вы введете некорректное выражение, то его статус будет обозначен как "Error in expression". Такое выражение вычисляться не будет.
Также, в консоли разработчика будет указан код ошибки.
![image](https://github.com/user-attachments/assets/5e780117-6192-4ba8-ad6e-8c1096d5b8ed)

После вычисления всех операций, статус выражения станет "Solved" и в поле result будет указано значение.
![image](https://github.com/user-attachments/assets/0bb2f2d1-529b-4fbf-9379-41af48d725bf)

Также, вы можете отслеживать какие задачи были созданы, статус их выполнения, если выберите All Tasks.
![image](https://github.com/user-attachments/assets/a0e7a27e-f37f-4d02-ab8f-1a6dbb9c5298)

Задачи назначаются оркестратором с задержкой в 5 секунд, так что для обновления просто снова нажимайте Receive.

Также вы можете обратиться к конкретной задаче или к конкретному выражению по их id. 
Для этого выберите One expression/One task и введите в появившееся поле id.

![image](https://github.com/user-attachments/assets/b70b1186-dbd1-408e-89d5-092e3e9edd1d)

## Конфигурирование
Все необходимые значения находятся в папке config в файле config.yaml. Вы можете менять их по своему усмотрению.
![image](https://github.com/user-attachments/assets/c1df1203-4e3e-4168-909c-2855f84b43f9)


## Ответ сервера
Сервер возвращает ответ в формате `{"expression":"2+3","result":5}"` для верного выражения и в формате `{"expression":"2++3","error":"internal server error: repeating operators"}"`

Сервер возвращает следующие статусы:
  - 200: При условии правильности переданного арифметического выражения, также будет возвращено полученное значение
  - 422: Если сервер не может обработать переданное выражение, например `{"expression": true}` или же в выражении содержатся неправильные символы, например `{"expression": "a+b"}`
  - 500: Если сервер смог обработать переданное выражение, но оно не имеет математического смысла, сюда входят ошибки деления на 0, неверные скобки, пустая строка, повторение операторов.

## Тестирование
К сожалению, автоматические тесты у меня сделать не получилось. Однако программа успешно выполняет асинхронные вычисления, в чем можно убедиться, если послать друг за другом несколько выражений.

Пример:
![image](https://github.com/user-attachments/assets/da10b230-2fc7-4114-bec7-4bfd342640b3)
![image](https://github.com/user-attachments/assets/92fa192a-99c7-4d0e-ab3c-8abf46e62d05)










