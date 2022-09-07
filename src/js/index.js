import Pusher from "pusher-js";
import { successAlert, errorAlert, warningAlert, Prompt } from "./attention";
import { FormValidation, FormSaveClose } from "./form_check";
import { monitoring } from "./monitoring_live";

let pusher = new Pusher(process.env.PUSHER_KEY, {
  authEndPoint: "/pusher/auth",
  wsHost: process.env.PUSHER_HOST,
  wsPort: process.env.PUSHER_PORT,
  forceTLS: false,
  enableTransports: ["ws", "wss"],
  disableTransports: [],
});

let publicChannel = pusher.subscribe("public-channel");

publicChannel.bind("app-starting", (data) => {
  successAlert(data.message);
});

publicChannel.bind("app-ending", (data) => {
  warningAlert(data.message);
});

export {
  pusher,
  publicChannel,
  FormValidation,
  FormSaveClose,
  successAlert,
  errorAlert,
  warningAlert,
  Prompt,
  monitoring,
};
