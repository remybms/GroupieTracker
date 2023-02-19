const artist = [
    {name:"ACDC"},
    {name:"SOJA"},
    {name:"Queen"},
    {name:"queez"},
    {name:"sonj"}
];

const searchinput = document.getElementById('search-header')
searchinput.AddEventListener('keyup', myResearch)

function myResearch(){
    const input = searchinput.value;
    console.log(input);
    const result = artist.filter(item => item.name.toLocaleLowerCase().includes(input.toLocaleLowerCase()));
    let final = '';
    if (input != ""){
        result.forEach(resultItem =>
                console.log(resultItem.name),
                final += `<div class="suggestion">${resultItem.name}</div>`
            )
    }
    document.getElementById("suggestion").innerHTML = final;
}
