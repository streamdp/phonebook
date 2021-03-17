let id_salt = 0;
let checkboxes_visible = false;
let arr_with_added_fields_Ids = [];
let phone_mask = "8 (999) 9-99-99";
let mobile_phone_mask = "+375 (99) 999-99-99";

$('#deleteModal').on('show.bs.modal', function (event) {
    const queryString = window.location.search;
    const urlParams = new URLSearchParams(queryString);
    const button = event.relatedTarget;
    const record_id = button.getAttribute('data-mdb-recordid');
    $('#deleteModalOkButton').attr("onclick", "location.href='/deleteRecord?id="+record_id+"&p="+urlParams.get('p')+"';");
});

$(function(){
    if (document.getElementById("inputGroupSelect03").value === "") {
        document.getElementById("third_select").hidden = true;
    } else {
        document.getElementById("third_select").hidden = false;
    }
});

$('#inputGroupSelect02').on('change', function () {
    let selected_option = document.getElementById('inputGroupSelect02').value;
    change_select(selected_option, "#inputGroupSelect03");
    document.getElementById("third_select").hidden = false;
});

$('#modalInputGroupSelect02').on('change', function () {
    let selected_option = document.getElementById('modalInputGroupSelect02').value;
    change_select(selected_option, "#modalInputGroupSelect03");
});

function del_fields(id) {
    $(id).remove();
    arr_with_added_fields_Ids.splice(arr_with_added_fields_Ids.indexOf(id), 1);
}

let checkboxes_should_be_removed = $('input[name="should_be_removed"]');
let array_checkboxes_ids = [];
let delete_button = $("#groupdeletedoaction");
let group_delete_button = $("#groupdelete");
function show_hide_checkboxes(){
   $.each(checkboxes_should_be_removed, function (index, item) {
        item.hidden = checkboxes_visible;
   })
    checkboxes_visible = !checkboxes_visible;
  if (checkboxes_visible) {
      delete_button.show();
      group_delete_button.toggleClass("btn-outline-light");
      group_delete_button.toggleClass("btn-outline-warning");
  } else {
      group_delete_button.toggleClass("btn-outline-light");
      group_delete_button.toggleClass("btn-outline-warning");
      delete_button.hide();
  }

}
checkboxes_should_be_removed.on('change', function (event){
    let checkbox = event.currentTarget;
    if ($(checkbox).is(":checked")){
        array_checkboxes_ids.push(checkbox.value);
    } else {
        array_checkboxes_ids.splice(array_checkboxes_ids.indexOf(checkbox.value), 1);
    }
});

function add_fields(id) {
    let val_input_fields = [];
    let fields = $("#" + id + " .form-control[name='" + id + "']");
    $.each(fields, function (index, item) {
        let a = [];
        a.push($(item).val());
        a.push($(item).attr('id'));
        val_input_fields.push(a);
    });
    let new_field_id = id + id_salt;
    document.getElementById(id).innerHTML += "<div id=\"div_"+new_field_id+"\" class=\"input-group mb-3 col-lg\"><input id=\"input_" + new_field_id + "\" name=\"" + id + "\" type=\"text\" class=\"form-control\" placeholder=\"Введите номер...\">" +
        "<div class=\"input-group-append\"><span class=\"btn btn-outline-danger\" type=\"button\" onClick=\"del_fields('#div_" + new_field_id + "')\">&minus;</span></div></div>"
    id_salt++;
    if (id === "service_mobile_num") {
        $("#input_" + new_field_id).mask(mobile_phone_mask);
    } else {
        $("#input_" + new_field_id).mask(phone_mask);
    }
    arr_with_added_fields_Ids.push("#div_" + new_field_id);
    for (let i = 0; i < val_input_fields.length; i++) {
        $("#" + val_input_fields[i][1]).val(val_input_fields[i][0]);
    }
    return "#input_" + new_field_id
}

function change_select(value, select_two) {
    $.getJSON("v1/third", {id: value}, function (response) {
        $(select_two + " option").remove();
        $(select_two).append(
            $("<option value=\"\" selected></option>").text("Выберите...")
        );
        $.each(response, function (index, item) {
            $(select_two).append(
                $("<option></option>").text(item).val(item)
            )
        });
    });
}

$(function(){
    $("#input_service_num").mask(phone_mask);
    $("#input_personal_num").mask(phone_mask);
    $("#input_service_mobile_num").mask(mobile_phone_mask);
})

let update_modal_div = $('#updateModal')
update_modal_div.on('hide.bs.modal', function (event) {
    if (arr_with_added_fields_Ids.length > 0) {
        while (arr_with_added_fields_Ids.length > 0) {
            del_fields(arr_with_added_fields_Ids[0]);
        }
    }
    $.each($("#updateForm .form-control"), function (index, item) {
        $(item).val("");
    })
});

function fill_input(response, field_name_in_response, field_name_in_html) {

    if (response[field_name_in_response] != null) {
        $.each(response[field_name_in_response], function (index, item) {
            if (index > 0) {
                let new_field_id = add_fields(field_name_in_html);
                let input = $(new_field_id);
                input.val(item);
                if (field_name_in_html === "service_mobile_num") {
                    input.mask(mobile_phone_mask);
                } else {
                    input.mask(phone_mask);
                }
            } else {
                let input = $("#input_" + field_name_in_html);
                input.val(item);
                if (field_name_in_html === "service_mobile_num") {
                    input.mask(mobile_phone_mask);
                } else {
                    input.mask(phone_mask);
                }
            }
        });
    }
}

update_modal_div.on('show.bs.modal', function (event) {
    const button = event.relatedTarget;
    let whatever = button.getAttribute('data-mdb-whatever') // Extract info from data-* attributes
    if (whatever !== "create"){
    $.getJSON("v1/get_record", {id: whatever}, function (response) {
        $('#modalInputGroupSelect01').val(response['department']['first_level']);
        $('#modalInputGroupSelect02').val(response['department']['second_level']);
        $("#modalInputGroupSelect02").change();
        window.setTimeout(function(){
            $('#modalInputGroupSelect03').val(response['department']['third_level']);
        }, 600);
        $(".form-control[name='modal_id_user']").val(response['id']);
        $(".form-control[name='modal_first_name']").val(response['first_name']);
        $(".form-control[name='modal_last_name']").val(response['last_name']);
        $(".form-control[name='modal_middle_name']").val(response['middle_name']);
        $(".form-control[name='modal_position']").val(response['position']);
        fill_input(response, 'service_number', 'service_num');
        fill_input(response, 'personal_number', 'personal_num');
        fill_input(response, 'service_mobile_number', 'service_mobile_num');
    });
    } else {
        $('#updatesubmitbutton').text("Создать");
        $('#updateModalLabel').text("Создать запись в телефонной книге:");
    }
})