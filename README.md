# Vigilate - Golang WebSockets 練習

## References

- [Working with WebSockets in Go (Golang)](https://www.udemy.com/course/working-with-websockets-in-go/)
- [Course code respiratory](https://github.com/tsawler/ws-udemy)
  - 已完成大部分的 UI (bootstrap)
  - 提供部份後端程式
- [IPÊ-An open source Pusher server](https://github.com/dimiro1/ipe)

## Goals

- 以 Go 建立一個即時監測網路服務的網路應用程式 (Web application)
- 程式需具備以下功能
  - 使用者
    - [x] 登入
    - [x] 登出
    - [x] 移除
    - [ ] 權限 (e.g. 超級使用者, 一般使用者)
  - 監測服務
    - [x] 加入新的監測服務
    - [x] 刪除監測服務
    - [x] 定時檢測服務是否還在線
    - [x] 紀錄每次檢測結果
    - [x] 檢測完成後即時更新結果
  - 可監測服務種類
    - [x] HTTP
    - [x] HTTPS
  - 關閉應用程式
    - [x] 優雅的關機(Graceful Shutdown)
