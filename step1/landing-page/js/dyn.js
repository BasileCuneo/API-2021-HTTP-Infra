$(function () {
	console.log("Hello world");
	function loadJson(){
		$.getJSON("/dyn/",function (data) {
			console.log(data);
			var message = "Nothing interesting";
			let n = Math.floor(Math.random() * data.length)
			if (data.length > 0) {
				message = data[n].firstname + " " + data[n].lastname;
				if(data[n].interest.length > 0){
					message += " is interested in " + data[n].interest;
				}if(data[n].skills.length > 0){
					message += ", has skills in " + data[n].skills;
				}
				if(data[n].interest.length == 0 && data[n].skills.length == 0){
					message += " has no interest and skills 😞";
				}
			}
			$("#dyntxt").text(message);
		});
	}
	setInterval(loadJson,5000);
});
