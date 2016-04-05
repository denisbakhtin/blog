$(document).ready(function(){
  //initialize dashboard sidebar slim scroll
  $('#sidebar-dashboard .slimscroll').slimScroll({
    height: '100%',
  });

  //initialize markdown editor
  $('#markdown-editor').markdownEditor({
    imageUpload: true,
    uploadPath: '/admin/upload',
  });

  //get markdown editor content
  $('#markdown-form').submit(function(event){
    //on submit get editor content
    $('#content').val($('#markdown-editor').markdownEditor('content'));
  });

});

