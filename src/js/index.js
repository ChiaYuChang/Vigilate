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

publicChannel.bind("host-service-status-change", (data) => {
  // successAlert(data.message);
  console.log("receive data from host-service-status-change channel");
  Prompt().toast({
    msg: data.message,
    icon: "info",
    timer: 30000,
    showCloseButton: true,
  });

  // update tables

  // host page
  // remove existing table row if it exist
  let exist = !!document.getElementById("host-service-" + data.host_service_id);
  if (exist) {
    let row = document.getElementById("host-service-" + data.host_service_id);
    row.parentNode.removeChild(row);
  }
  // update tables, if t
  let tableExists = !!document.getElementById(data.status + "-table");
  if (tableExists) {
    console.log("add row!!");
    let tableRef = document.getElementById(data.status + "-table");
    let newRow = tableRef.tBodies[0].insertRow(-1);

    newRow.setAttribute("id", "host-service-" + data.host_service_id);
    let frstCell = newRow.insertCell(0);

    if (data.status === "healthy") {
      frstCell.innerHTML = `<span>${data.icon}</span> ${data.service_name}`;
    } else {
      frstCell.innerHTML = `<span>${data.icon}</span> ${data.service_name}
      <span class="badge bg-secondary pointer"
          onclick="checkNow(${data.host_service_id}, "${data.status}")">
          Click Now
      </span>`;
    }

    let scndCell = newRow.insertCell(1);
    if (data.status !== "pending") {
      scndCell.innerHTML = `${data.last_check}`;
    } else {
      scndCell.innerHTML = "Pending";
    }
    let thrdCell = newRow.insertCell(2);
  }
});

publicChannel.bind("host-service-count-change", (data) => {
  let healthyCountExists = !!document.getElementById("healthy_count");
  if (healthyCountExists) {
    if (data.healthy_count === "1") {
      document.getElementById("healthy_count").innerHTML =
        data.healthy_count + " Healthy service";
    } else {
      document.getElementById("healthy_count").innerHTML =
        data.healthy_count + " Healthy services";
    }

    if (data.warning_count === "1") {
      document.getElementById("warning_count").innerHTML =
        data.warning_count + " Warning service";
    } else {
      document.getElementById("warning_count").innerHTML =
        data.warning_count + " Warning services";
    }

    if (data.problem_count === "1") {
      document.getElementById("problem_count").innerHTML =
        data.problem_count + " Problem service";
    } else {
      document.getElementById("problem_count").innerHTML =
        data.problem_count + " Problem services";
    }

    if (data.pending_count === "1") {
      document.getElementById("pending_count").innerHTML =
        data.pending_count + " Pending service";
    } else {
      document.getElementById("pending_count").innerHTML =
        data.pending_count + " Pending services";
    }
  }
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
