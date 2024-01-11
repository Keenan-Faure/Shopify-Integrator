import {useEffect} from 'react';
import {useState} from "react";
import $ from 'jquery';
import '../../CSS/login.css';

function Export_Product()
{
    const ExportProduct = (event) =>
    {
        event.preventDefault();
        
        const api_key = localStorage.getItem('api_key');
        $.ajaxSetup({ headers: { 'Authorization': 'ApiKey ' + api_key} });
        $.get("http://localhost:8080/api/products/export", [], [], 'json')
        .done(function( _data) 
        {
            console.log(_data.message);

            
            let a_tag = document.createElement("a");
            a_tag.className = "tablink";
            a_tag.setAttribute("href", _data.message);
            a_tag.setAttribute("target", "_blank");
            a_tag.click();
            
        })
        .fail( function(xhr) 
        {
            alert(xhr.responseText);
        });
        
    }

    useEffect(() =>
    {

        /* Fix any incorrect elements */
        let navigation = document.getElementById("navbar");
        let modal = document.getElementById("model");
        modal.style.display = "block";
        navigation.style.animation = "MoveRight 1.2s ease";
        navigation.style.position = "fixed";
        navigation.style.left = "0%";
        navigation.style.width = "100%";

    }, []);

    return (
        <>
            <div className = 'modal1' id = "model" style={{background:'linear-gradient(to bottom, #202020c7, #111119f0)', zIndex: '2'}}>

                <form className = 'modal-content' style ={{backgroundColor: 'none'}} method = 'post' onSubmit={(event) => ExportProduct(event)} autoComplete='off' id = 'form1'>
                    <div className = 'modal-container'style = {{backgroundColor: 'transparent'}} id = "main">
                        <label style = {{fontSize: '18px', textDecoration:'underline'}}><b>Export Product</b></label>
                        <br /><br /><br />
                        <label><b>Exports existing products to a CSV File</b></label>
                        <br /><br /><br />
                        <button className = 'button' type = 'submit'>Export</button>
                    </div>
                </form>
            </div>    
        </>
    );
};
export default Export_Product;