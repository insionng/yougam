$(document).ready(function(){
	
	$('.data-table').dataTable({
		"bJQueryUI": true,
		"sPaginationType": "full_numbers",
		"sDom": '<""l>t<"F"fp>'
	});
	
	$('input[type=checkbox],input[type=radio],input[type=file]').uniform();
	
	$('select').select2();
	
	$("span.icon input:checkbox, th input:checkbox").click(function() {
		var checkedStatus = this.checked;
		var checkbox = $(this).parents('.widget-box').find('tr td:first-child input:checkbox');		
		checkbox.each(function() {
			this.checked = checkedStatus;
			if (checkedStatus == this.checked) {
				$(this).closest('.checker > span').removeClass('checked');
			}
			if (this.checked) {
				$(this).closest('.checker > span').addClass('checked');
			}
		});
	});

	$("#delrows").click(function() {

		$('#delrowids').attr('value',"");

	    var oTable = $('.data-table').dataTable();
	    $('input[type=checkbox]:checked', oTable.fnGetNodes()).each( function() {
	    		var v = $('#delrowids').attr('value');

				if (v != "") {
					v = v + "," + $(this).attr('data');
				}else{
					v = $(this).attr('data');
				};
		       $('#delrowids').attr('value',v);

 		});
 		
 		if ($('#delrowids').attr('value') != "") {$('#iform').submit();};

	});
});