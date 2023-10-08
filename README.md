# Go WebSockets 練習

- [Working with WebSockets in Go (Golang)](https://www.udemy.com/course/working-with-websockets-in-go/)
- [Course code respiratory](https://github.com/tsawler/ws-udemy)
  - 已完成大部分的 UI (bootstrap)
  - 提供部份後端程式

# 目標

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

# 實踐

- 以 PostgreSQL 作為資料庫
- 以 WebSocket 作為連線方式
- 以 webpack 打包 javascript 檔案
-
