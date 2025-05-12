function Prompt() {
    let toast = function(c) {
      console.log("Toast");

      const {
        msg= "",
        icon = "success",
        position = "top-end",
      } = c;

      const Toast = Swal.mixin({
        toast: true,
        title: msg,
        position: position,
        showConfirmButton: false,
        timer: 3000,
        icon: icon,
        timerProgressBar: true,
        didOpen: (toast) => {
          toast.onmouseenter = Swal.stopTimer;
          toast.onmouseleave = Swal.resumeTimer;
        }
      });
      Toast.fire({});
    }

    let success = function(c) {
      console.log("Alert");
      const {
        title = "",
        htmlText = "",
        icon = "",
        footer = ""
      } = c;

      Swal.fire({
        icon: icon,
        title: title,
        text: htmlText,
        footer: footer
      });

    }

    let error = function(c) {
      console.log("error");
      const {
        title = "",
        htmlText = "",
        icon = "",
        footer = ""
      } = c;

      Swal.fire({
        icon: icon,
        title: title,
        text: htmlText,
        footer: footer
      });

    }

    async function custom(c) {
      const {
        icon = "",
        msg = "",
        title = "",
        showConfirmButton = true,
      } = c;

      const { value: result } = await Swal.fire({
        icon: icon,
        title: title,
        html: msg,
        focusConfirm: false,
        backdrop: false,
        showCancelButton: true,
        showConfirmButton: showConfirmButton,
        willOpen: () =>{
          if (c.willOpen !== undefined) {
            c.willOpen();
          }
        },
        preConfirm: () => {
          return [
            document.getElementById("start").value,
            document.getElementById("end").value
          ];
        },
        didOpen: () => {
          if (c.didOpen !== undefined) {
            c.didOpen();
          }
        }
      });


      if (result) {
        if (result.dismiss !== Swal.DismissReason.cancel) {
            if (result.value !== "") {
                if (c.callback !== undefined ) {
                    c.callback(result);
                } else {
                    c.callback(false);
                }
            }
        } else {
            c.callback(false);
        }
      } 
    }

    return {
      toast: toast,
      success: success,
      error: error,
      custom: custom,
    }
  }