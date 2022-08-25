import Pusher from "pusher-js";
import { successAlert, errorAlert, warningAlert, Prompt } from "./attention";
import { FormValidation, FormSaveClose } from "./form_check";

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

// window.successAlert = successAlert;
// window.errorAlert = errorAlert;
// window.warningAlert = warningAlert;
// window.Prompt = Prompt;
// window.FormValidation = FormValidation;
// window.FormSaveClose = FormSaveClose;
export {
  pusher,
  publicChannel,
  FormValidation,
  FormSaveClose,
  successAlert,
  errorAlert,
  warningAlert,
  Prompt,
};
