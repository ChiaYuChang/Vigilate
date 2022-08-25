import { errorAlert } from "./attention";

export { FormValidation, FormSaveClose };

function FormValidation(id, self) {
  const val = () => {
    document.getElementById("action").value = 0;
    //   let form = document.getElementById("host-form");
    let form = document.getElementById(id);
    if (form.checkValidity() === false) {
      errorAlert("Error: check all tabs!");
      self.event.preventDefault();
      self.event.stopPropagation();
    }
    form.classList.add("was-validated");

    if (form.checkValidity() === true) {
      form.submit();
    }
  };
  return val;
}

function FormSaveClose(id, self) {
  const saveClose = () => {
    document.getElementById("action").value = 1;
    let form = document.getElementById(id);
    if (form.checkValidity() === false) {
      errorAlert("Error: check all tabs!");
      self.event.preventDefault();
      self.event.stopPropagation();
    }
    form.classList.add("was-validated");

    if (form.checkValidity() === true) {
      form.submit();
    }
  };
  return saveClose;
}
