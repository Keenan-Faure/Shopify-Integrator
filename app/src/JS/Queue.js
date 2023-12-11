import { createRoot } from 'react-dom/client';
import { flushSync } from 'react-dom';
import {useEffect, useState} from 'react';
import $ from 'jquery';
import Page1 from '../components/Page1';
import Queue_details from '../components/semi-components/queue-details';


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
            let root;
            let pan_main = document.querySelector(".pan-main");
            let div = document.createElement("div");
            pan_main.appendChild(div);

            root = createRoot(div);
            root.render(_data.map((el, i) => <Queue_details key={`${el.title}_${i}`}
            />))
            
        })
        .fail( function(xhr) { alert(xhr.responseText); });


    }, []);
    
    return (
        <div className = "queue">
            <div className = "main">
                <div className = "search">
                    <form className = "search-area" id = "search" autoComplete='off' onSubmit={(event) => SearchProduct(event)}>
                        <input className ="search-area" type="search" placeholder="Search..." 
                        name = "search" value = {inputs.search || ""}  onChange = {handleChange}></input>
                    </form>    
                </div>
                <div className = "main-elements">
                    <div className = "empty-message">No results found.</div>
                    <div className = "pan-main" id = "pan-main">

                    </div>
                </div>
                <div className = "center" id = "pag"></div>
            </div>

            <Page1 title = "Products"/>
            <div className = "details">
                <div className = 'close-button'>&times;</div>
            </div>

        </div>
    );
}

export default Queue;