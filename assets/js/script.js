jQuery(document).ready(function(){
	jQuery(".tab a").click(function(){
		jQuery("form input[name='target']").val(jQuery(this).data("target"));
		jQuery("form").submit();
	});
});