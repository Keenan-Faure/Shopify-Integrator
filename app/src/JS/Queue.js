import { createRoot } from 'react-dom/client';
import { flushSync } from 'react-dom';
import {useEffect, useState} from 'react';
import $ from 'jquery';
import Page2 from '../components/Page2';
import Queue_details from '../components/semi-components/queue-details';
import Detailed_queue from '../components/semi-components/Queue/detailed_queue';


import '../CSS/page1.css';

function Queue()
{
    const[inputs, setInputs] = useState({});

    const handleChange = (event) =>
    {
        const name = event.target.name;
        const value = event.target.value;
        setInputs(values => ({...values, [name]: value}))
    }

    const SearchProduct = (event) =>
    {
        event.preventDefault();

    }

    useEffect(()=> 
    {
        /* Ensures the page elements are set correctly */
        let navigation = document.getElementById("navbar");
        window.onload = function(event)
        {
            navigation.style.left = "30%";
            navigation.style.position = "absolute";
            navigation.style.width = "70%";
            navigation.style.animation = "MoveLeft 0.8s ease";
        }

        /*  API INITIAL-REQUEST */
        const api_key = localStorage.getItem('api_key');
        $.ajaxSetup({ headers: { 'Authorization': 'ApiKey ' + api_key} });
        $.get("http://localhost:8080/api/queue?page=1", [], [])
        .done(function( _data) 
        {
            console.log(_data);
            if(_data == "")
            {
                document.querySelector(".empty-message").style.display = "block";
            }
            else 
            {
                let root;
                let pan_main = document.querySelector(".pan-main");
                let div = document.createElement("div");
                pan_main.appendChild(div);

                root = createRoot(div);
                root.render(_data.map((el, i) => <Queue_details key={`${el.title}_${i}`}
                Queue_Type={el.queue_type} Queue_Instruction={el.instruction} Queue_Status={el.status} Queue_ID={el.id}
                />))
            } 
        })
        .fail( function(xhr) { alert(xhr.responseText); });

        /* When the user clicks on the pan elements show info about that specified pan element */
        function DetailedView()
        {
            let queue = document.querySelector(".queue");
            let pan = document.querySelectorAll(".pan");
            for(let i = 0; i < pan.length; i++)
            {
                pan[i].addEventListener("click", () =>
                {
                    let id = pan[i].querySelector(".p-d-id").innerHTML;
                    /*  API  */
                    const api_key = localStorage.getItem('api_key');
                    $.ajaxSetup({ headers: { 'Authorization': 'ApiKey ' + api_key} });
                    $.get("http://localhost:8080/api/queue/" + id, [], [], 'json')
                    .done(function(_data) 
                    {   
                        console.log(_data);
                        if(document.querySelector(".details") != null)
                        //div already exists, remove it, and create another
                        {

                            document.querySelector(".details").remove();
                            let details = document.createElement('div');
                            details.className = "details";
                            queue.appendChild(details);

                            let rot = createRoot(details);
                            rot.render( <Detailed_queue key={`${_data.title}_${i}`} Queue_Status={_data.status} Queue_Description={_data.description}
                            Queue_Type={_data.queue_type} Queue_Instruction={_data.instruction} Queue_ID={_data.id} Created_At={_data.created_at} 
                            Updated_At={_data.updated_at}
                            />)
                        }
                        else 
                        //create new div
                        {
                            let details = document.createElement('details');
                            queue.appendChild(details);
                            let rot = createRoot(details);
                            rot.render( <Detailed_queue key={`${_data.title}_${i}`} Queue_Status={_data.status} Queue_Description={_data.description}
                            Queue_Type={_data.queue_type} Queue_Instruction={_data.instruction} Queue_ID={_data.id} Created_At={_data.created_at} 
                            Updated_At={_data.updated_at}
                            />)
                        }
                    })
                    .fail( function(xhr) 
                    {
                        alert(xhr.responseText);
                    });
                    setTimeout(() =>
                    {
                        let filter = document.querySelector(".filter");
                        let main = document.querySelector(".main");
                        let navbar = document.getElementById("navbar");
                        let details = document.querySelector(".details");

                        filter.style.animation = "Fadeout 0.5s ease-out";
                        main.style.animation = "Fadeout 0.5s ease-out";
                        navbar.style.animation = "Fadeout 0.5s ease-out";
                        filter.style.display = "none";
                        main.style.display = "none";
                        navbar.style.display = "none";
                        details.style.display = "block";
                    }, 50);
                });
            } 
        }

        /* Script to automatically format the number of elements on each page */
        const content = document.querySelector('.center'); 
        const paginationContainer = document.createElement('div');
        const paginationDiv = document.body.appendChild(paginationContainer);
        paginationContainer.classList.add('pagination');
        content.appendChild(paginationContainer);

        document.querySelector(".pan-main").remove();
        let div = document.createElement("div");
        div.className = "pan-main";
        div.id = "pan-main";
        let main = document.querySelector(".main-elements");
        main.appendChild(div);

        function Pagintation(index)
        {
            let ahead = index + 1;
            /*  API  */
            $.get('http://localhost:8080/api/queue?page=' + ahead, [], [])
            .done(function( _data) 
            {
                console.log(_data);
                if(_data == "")
                {
                    let next = document.getElementById("next");
                    next.style.cursor = "not-allowed";
                    next.disabled = true;
                } 
            })
            .fail( function(xhr) { alert(xhr.responseText); });

            /* Check done to remove old elements if they exist */
            if(document.getElementById("next") != null && document.getElementById("prev") != null && document.getElementById("hod") != null)
            //If they exist remove them, and create new based on the new index value
            {
                document.getElementById("next").remove();
                document.getElementById("prev").remove();
                document.getElementById("hod").remove();

                const pageButton = document.createElement('button');
                pageButton.id = "hod";
                pageButton.className = "active";
                pageButton.innerHTML = index;
                paginationDiv.appendChild(pageButton);

                const nextPage = document.createElement('button');
                nextPage.id = "next";
                nextPage.innerHTML = "→";
                paginationDiv.appendChild(nextPage);

                const prevPage = document.createElement('button');
                prevPage.id = "prev";
                prevPage.innerHTML = "←";
                paginationDiv.appendChild(prevPage);
                if(index == 1) { prevPage.disabled = true; prevPage.style.cursor = "not-allowed"; }
                else if(index > 1) { prevPage.style.cursor = "pointer"; prevPage.disabled = false; nextPage.disabled = false; }
                else if(index <= 1) {prevPage.disabled = true; prevPage.style.cursor = "not-allowed"; }
                

                nextPage.addEventListener("click", () =>
                {
                    index = index + 1;
                    /* Fetches the data from page, based on the page / index value */
                    const page = "http://localhost:8080/api/queue?page=" + index;
                    /*  API  */
                    const api_key = localStorage.getItem('api_key');
                    $.ajaxSetup({ headers: { 'Authorization': 'ApiKey ' + api_key} });
                    $.get(page, [], [])
                    .done(function( _data) 
                    {
                        console.log(_data);
                        document.querySelector(".pan-main").remove();
                        let div = document.createElement("div");
                        div.className = "pan-main";
                        div.id = "pan-main";
                        let main = document.querySelector(".main-elements");
                        main.appendChild(div);
                        let root = createRoot(div);
                        flushSync(() => 
                        { 
                            root.render(_data.map((el, i) => <Queue_details key={`${el.title}_${i}`}
                            Queue_Type={el.queue_type} Queue_Instruction={el.instruction} Queue_Status={el.status} Queue_ID={el.id} 
                            />))
                        });
                    })
                    .fail( function(xhr) { alert(xhr.responseText); });

                    let ahead = index + 1;
                    /*  API  */
                    $.get('http://localhost:8080/api/products?page=' + ahead, [], [])
                    .done(function( _data) 
                    {
                        console.log(_data);
                        if(_data == "")
                        {
                            let next = document.getElementById("next");
                            next.style.cursor = "not-allowed";
                            next.disabled = true;
                        } 
                    })
                    .fail( function(xhr) { alert(xhr.responseText); });

                    Pagintation(index);
                    setTimeout(() => { DetailedView();}, 200);
                });

                prevPage.addEventListener("click", () =>
                {
                    index = index - 1;
                    /* Fetches the data from page, based on the page / index value */
                    const page = "http://localhost:8080/api/queue?page=" + index;

                    /*  API  */
                    const api_key = localStorage.getItem('api_key');
                    $.ajaxSetup({ headers: { 'Authorization': 'ApiKey ' + api_key} });
                    $.get(page, [], [])
                    .done(function( _data) 
                    {
                        console.log(_data);

                        document.querySelector(".pan-main").remove();
                        let div = document.createElement("div");
                        div.className = "pan-main";
                        div.id = "pan-main";
                        let main = document.querySelector(".main-elements");
                        main.appendChild(div);
                        let root = createRoot(div);

                        flushSync(() => 
                        { 
                            root.render(_data.map((el, i) => <Queue_details key={`${el.title}_${i}`} 
                            Queue_Type={el.queue_type} Queue_Instruction={el.instruction} Queue_Status={el.status} Queue_ID={el.id}
                            />))
                        });
                    })
                    .fail( function(xhr) { alert(xhr.responseText); });
                    Pagintation(index);
                    setTimeout(() => { DetailedView();}, 200);
                });
            }
            else 
            //If they dont exist create new ones 
            {
                const pageButton = document.createElement('button');
                pageButton.id = "hod";
                pageButton.className = "active";
                pageButton.innerHTML = index;
                paginationDiv.appendChild(pageButton);

                const nextPage = document.createElement('button');
                nextPage.id = "next";
                nextPage.innerHTML = "→";
                paginationDiv.appendChild(nextPage);

                const prevPage = document.createElement('button');
                prevPage.id = "prev";
                prevPage.innerHTML = "←";
                paginationDiv.appendChild(prevPage);

                if(index == 1) { prevPage.disabled = true; prevPage.style.cursor = "not-allowed"; }
                else if(index > 1) { prevPage.style.cursor = "pointer"; prevPage.disabled = false; nextPage.disabled = false; }
                else if(index <= 1) {prevPage.disabled = true; prevPage.style.cursor = "not-allowed"; }
                nextPage.addEventListener("click", () =>
                {
                    index = index + 1;
                    /* Fetches the data from page, based on the page / index value */
                    const page = "http://localhost:8080/api/queue?page=" + index;
                    /*  API  */
                    const api_key = localStorage.getItem('api_key');
                    $.ajaxSetup({ headers: { 'Authorization': 'ApiKey ' + api_key} });
                    $.get(page, [], [])
                    .done(function( _data) 
                    {
                        console.log(_data);
                        document.querySelector(".pan-main").remove();
                        let div = document.createElement("div");
                        div.className = "pan-main";
                        div.id = "pan-main";
                        let main = document.querySelector(".main-elements");
                        main.appendChild(div);
                        let root = createRoot(div);

                        flushSync(() => 
                        { 
                            root.render(_data.map((el, i) => <Queue_details key={`${el.title}_${i}`}
                            Queue_Type={el.queue_type} Queue_Instruction={el.instruction} Queue_Status={el.status} Queue_ID={el.id}
                            />))
                        });
                    })
                    .fail( function(xhr) { alert(xhr.responseText); });
                    Pagintation(index);
                    setTimeout(() => { DetailedView();}, 200);
                });
                prevPage.addEventListener("click", () =>
                {
                    index = index - 1;
                    /* Fetches the data from page, based on the page / index value */
                    const page = "http://localhost:8080/api/queue?page=" + index;

                    /*  API  */
                    const api_key = localStorage.getItem('api_key');
                    $.ajaxSetup({ headers: { 'Authorization': 'ApiKey ' + api_key} });
                    $.get(page, [], [])
                    .done(function( _data) 
                    {
                        console.log(_data);
                        document.querySelector(".pan-main").remove();
                        let div = document.createElement("div");
                        div.className = "pan-main";
                        div.id = "pan-main";
                        let main = document.querySelector(".main-elements");
                        main.appendChild(div);
                        let root = createRoot(div);

                        flushSync(() => 
                        { 
                            root.render(_data.map((el, i) => <Queue_details key={`${el.title}_${i}`}
                            Queue_Type={el.queue_type} Queue_Instruction={el.instruction} Queue_Status={el.status} Queue_ID={el.id}
                            />))
                        });
                    })
                    .fail( function(xhr) { alert(xhr.responseText); });
                    Pagintation(index);
                    setTimeout(() => { DetailedView();}, 200);
                });
            } 
        }
        Pagintation(1);
        setTimeout(() => { DetailedView();}, 200);

        let C_filter = document.getElementById("clear_filter");
        C_filter.addEventListener("click", () => 
        {
            document.querySelector(".empty-message").style.display = "none";

            const api_key = localStorage.getItem('api_key');
            $.ajaxSetup({ headers: { 'Authorization': 'ApiKey ' + api_key} });
            $.get("http://localhost:8080/api/queue?page=1", [], [])
            .done(function( _data) 
            {
                console.log(_data);
                let filter_button = document.getElementById("_filter");
                let C_filter = document.getElementById("clear_filter");
                filter_button.disabled = true;
                C_filter.disabled = true;
                filter_button.style.cursor = "not-allowed";
                C_filter.style.cursor = "not-allowed";

                let root;
                let pan_main;
                if(document.querySelector(".pan-main") != null){ document.querySelector(".pan-main").remove(); }
            
                pan_main = document.createElement('div');
                let main_elements = document.querySelector(".main-elements");
                pan_main.className = "pan-main";
                main_elements.appendChild(pan_main);
                root = createRoot(pan_main);
                root.render(_data.map((el, i) => <Queue_details key={`${el.title}_${i}`} Queue_Updated_At={el.updated_at} Queue_Creation_Date={el.created_at} 
                Queue_Type={el.queue_type} Queue_Instruction={el.instruction} Queue_Status={el.status} Queue_ID={el.id}
                />))
                setTimeout(() => { DetailedView();}, 200);
                Pagintation(1);
            })
            .fail( function(xhr) { alert(xhr.responseText); });
        });


    }, []);
    
    return (
        <div className = "queue">
            <div className = "main">
                <div className = "search">
                    <form className = "search-area" style ={{top: '30px'}} id = "search" autoComplete='off' onSubmit={(event) => SearchProduct(event)}>
                        <input className ="search-area" type="search" placeholder="Search..." 
                        name = "search" value = {inputs.search || ""}  onChange = {handleChange}></input>
                    </form>    
                </div>
                <div className = "main-elements" style={{top: '53%'}}>
                    <div className = "empty-message">No results found.</div>
                    <div className = "pan-main" id = "pan-main">

                    </div>
                </div>
                <div className = "center" id = "pag"></div>
            </div>

            <Page2 title = "Queue"/>
            <div className = "details">
                <div className = 'close-button'>&times;</div>
            </div>
        </div>
    );
}

export default Queue;