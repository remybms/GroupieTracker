const artist = [
    {name:"ACDC"},
    {name:"SOJA"},
    {name:"Queen"},
];

const searchinput = document.getElementById('searchInput')

searchInput.AddEventListener('keyup', function(){
    const input = searchinput.value;

    const result = artist.filter(item => item.name.toLocaleLowerCase().includes(input.toLocaleLowerCase()));

    let suggestion = '';

    if (input != ""){
    result.forEach(resultItem => 
            suggestion +=`
            <div class="suggestion">${resultItem.name}</div>
            `
        )
    }
    
    document.getElementById().innerHTML = suggestion;
})
