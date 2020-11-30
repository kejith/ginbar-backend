function loadPosts(url) {
    $.getJSON(url, function(data) {
        console.log(data);

        data.forEach(element => {
            var postElement = '<div id="'+ containerID +'" class="col-md-2 post-container"> \
            <a href="'+ anchorHref +'" class="post-thumbnail img-fluid"> \
            <img alt="picture" src="'+ imageSrc +'" /></a></div>';
            
            var content = document.getElementById("content").appendChild(postElement);
        });
    });
}

function LinkHandler() {
    return false;
}



$(document).ready(function(){
    var url = "http://127.0.0.1:8080/api/post";

    loadPosts(url);

    $(".alert").hide();

    $("#upload-form").submit(function(event){
        event.preventDefault();

        $.ajax({
            type: "POST",
            url: "api/post/create",
            data: $("form#upload-form").serialize(),
            dataType: "json",
            encode: true,
            success: function(){
                console.log("Upload successfull");
                $("#upload-alert-success").show();
            },
            error: function(XMLHttpRequest, textStatus, errorThrown) {
                console.log("Upload Form => Status: " + textStatus + " Error: " + errorThrown);
            }
        });
    });

    
})


