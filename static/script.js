function loadMail(id, index) {
    if (typeof index !== 'undefined') {
        window.open('/mail/multi/'+id+'/'+index, '_blank');
    }
    else {
        window.open('/mail/single/'+id, '_blank');
    }
}

function loadMailMeta(id) {
    var xmlhttp;

    xmlhttp = new XMLHttpRequest();

    xmlhttp.onreadystatechange = function() {
        if (xmlhttp.readyState == XMLHttpRequest.DONE ) {
           if(xmlhttp.status == 200){
               document.getElementById("meta").innerHTML = xmlhttp.responseText;
           }
           else if(xmlhttp.status == 400) {
              alert('There was an error 400')
           }
           else {
               alert('something else other than 200 was returned')
           }
        }
    };

    xmlhttp.open("GET", "/mail/meta/" + id, true);
    xmlhttp.send();
}
