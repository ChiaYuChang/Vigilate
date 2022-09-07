export { monitoring };
import { Prompt } from "./attention";

function monitoring(element) {
  if (!element.checked) {
    Prompt().confirm({
      html: "This will stop monitoring of all hosts and service, Are you sure?",
      callback: function (result) {
        if (result) {
          // want to turn monitoring off
          console.log("Turning monitoring off");
        } else {
          console.log("Canceling the operation");
          element.checked = true;
        }
      },
    });
  } else {
    console.log("It's off!!!");
  }
}
