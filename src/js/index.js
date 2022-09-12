import Pusher from "pusher-js";
import { successAlert, errorAlert, warningAlert, Prompt } from "./attention";
import { FormValidation, FormSaveClose } from "./form_check";
import { monitoring } from "./monitoring_live";

function capitalizeFirstLetter(string) {
  return string.charAt(0).toUpperCase() + string.slice(1);
}

function removeRowInServerStatusOverviewTable(data) {
  let tableRef = document.getElementById(
    "overview-" + data.old_status + "-srv"
  );

  let row = document.getElementById("overview-hs-" + data.host_service_id);
  if (tableRef.tBodies[0].childElementCount === 1) {
    row.setAttribute("id", "no-services");
    row.innerHTML = `<td colspan="5">No services</td>`;
  } else {
    row.parentNode.removeChild(row);
  }
}

function addRowInServerStatusOverviewTable(data) {
  if (!!document.getElementById("no-services")) {
    let row = document.getElementById("no-services");
    row.parentNode.removeChild(row);
  }

  let tableRef = document.getElementById("overview-" + data.status + "-srv");
  let newRow = tableRef.tBodies[0].insertRow(-1);
  let bg_type = "bg-success";
  switch (data.status) {
    case "warning":
      bg_type = "bg-warning";
    case "pending":
      bg_type = "bg-secondary";
    case "problem":
      bg_type = "bg-danger";
  }
  newRow.setAttribute("id", "overview-hs-" + data.host_service_id);
  newRow.innerHTML = `
    <td><a href="/admin/host/${data.host_id}#problem-content">${
    data.host_name
  }</a></td>
    <td>${data.service_name}</td>
    <td>
        <span class="badge ${bg_type}">${capitalizeFirstLetter(
    data.status
  )}</span>
    </td>
    <td>""</td>
    `;
}

function updateServerStatusTable(data) {
  // console.log("add row!!");
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

  if (!!document.getElementById("schedule-table")) {
    let tableRef = document.getElementById("schedule-table");
    tableRef.tBodies[0].innerHTML = "";

    let row = tableRef.tBodies[0].insertRow(-1);
    row.setAttribute("id", "no-scheduled-checks");
    row.innerHTML = `<td colspan="5">No scheduled checks!</td>`;
  }
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

  // update service status tables, if server status has changed
  if (!!document.getElementById(data.status + "-table")) {
    updateServerStatusTable(data);
  }

  // update service status overview table, if server status has changed
  // remove row if client is browsing old status page
  if (!!document.getElementById("overview-" + data.old_status + "-srv")) {
    console.log("remove one row: " + "overview-hs-" + data.host_service_id);
    removeRowInServerStatusOverviewTable(data);
  }
  // add row if client is browsing new status page
  if (!!document.getElementById("overview-" + data.status + "-srv")) {
    console.log("add one row");
    addRowInServerStatusOverviewTable(data);
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

publicChannel.bind("schedule-changed-event", (data) => {
  console.log("schedule change event");
  if (!!document.getElementById("schedule-table")) {
    let tableRef = document.getElementById("schedule-table");

    if (!!document.getElementById("no-scheduled-checks")) {
      let row = document.getElementById("no-scheduled-checks");
      row.parentNode.removeChild(row);
    }

    let row_content = [
      data.host_name,
      data.service_name,
      data.schedule,
      data.last_run,
      data.next_run,
    ];
    if (!!document.getElementById(`schedule-${data.host_service_id}`)) {
      // if the target service is already in the table
      // console.log("find target row");
      let row = document.getElementById(`schedule-${data.host_service_id}`);
      for (let i = 0; i < row_content.length; i++) {
        row.cells[i].innerHTML = row_content[i];
      }
    } else {
      // if the target service is missing
      // console.log("cannot find target row");
      let row = tableRef.tBodies[0].insertRow(-1);
      row.setAttribute("id", `schedule-${data.host_service_id}`);
      for (let i = 0; i < row_content.length; i++) {
        let cell = row.insertCell(i);
        let text = document.createTextNode(row_content[i]);
        cell.appendChild(text);
        console.log(row_content[i]);
      }
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
