function loadMailMeta(e,self) {
    e.preventDefault();
    document.getElementById("meta").data = self.href;
    return;

    var xmlhttp = new XMLHttpRequest();

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

    xmlhttp.open("GET", self.href, true);
    xmlhttp.send();
}
