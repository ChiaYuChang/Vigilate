import Pusher from "pusher-js";
import { successAlert } from "./attention";

let pusher = new Pusher(process.env.PUSHER_KEY, {
  authEndPoint: "/pusher/auth",
  wsHost: process.env.PUSHER_HOST,
  wsPort: process.env.PUSHER_PORT,
  forceTLS: false,
  enableTransports: ["ws", "wss"],
  disableTransports: [],
});

let publicChannel = pusher.subscribe("public-channel");

publicChannel.bind("test-event", (data) => {
  successAlert(data.message);
});
