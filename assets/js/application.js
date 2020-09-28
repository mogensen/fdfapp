require("expose-loader?$!expose-loader?jQuery!jquery");
require("bootstrap/dist/js/bootstrap.bundle.js");

$(() => {

    $("#activity-Duration").on( "change", function() { 
      if ($(this).val() > 9) {
        $("#activity-duration-warning").show()
      } else {
        $("#activity-duration-warning").hide()
      }
    }); 

    $('input[name="participants"]').on( "change", function() { 
      checked = $('input[name="participants"]:checked').length
      if (checked < 4) {
        $("#activity-participants-warning").show()
      } else {
        $("#activity-participants-warning").hide()
      }
    }); 
});
