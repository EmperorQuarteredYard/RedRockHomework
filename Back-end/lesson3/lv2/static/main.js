document.getElementById("submitData").addEventListener("click",(event)=>{event.preventDefault();PostData()})

function PostData(){
    let scoreForm = document.getElementById("scores")
    let scores=scoreForm.getElementsByTagName("input")
    let studentScore=[]
    for (let i = 0; i < scores.length; i++) {
        studentScore.push(parseFloat(scores[i].value))
        console.log(typeof studentScore[0])
    }
    let postData={
        name:(document.getElementById("nameInput").value.trim()),
        score:studentScore
    }
    console.log(JSON.stringify(postData))
    fetch("/calculate",{
        method:"POST",
        headers:{
            'Content-Type':"application/json"
        },
        body:JSON.stringify(postData)
    })
        .then(response=>{
            if(!response.ok){
                console.table(response)
                throw new Error('服务器响应错误: ' + response.status)
            }
            return response.json()
        })
        .then(data=>{
            document.getElementById("resultContent").textContent=data.average;
            document.getElementById('result').style.display = 'block';

            // 滚动到结果位置
            document.getElementById('result').scrollIntoView({ behavior: 'smooth' });
        })
}