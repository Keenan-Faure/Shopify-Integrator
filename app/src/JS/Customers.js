import {useEffect, useState} from 'react';
import $ from 'jquery';
import Customer_details from '../components/semi-components/customer-details';
import Page1 from '../components/Page1';
import '../CSS/page1.css';

/* Must start with a Caps letter */
function Customers()
{
    const [data, setData] = useState([]);

    useEffect(()=> 
    {
        /* Ensures the navbar is set correctly */
        let navigation = document.getElementById("navbar");
        let model = document.getElementById("model");
        let main = document.querySelector(".main");
        window.onload = function(event)
        {
            navigation.style.left = "0%";
            navigation.style.position = "relative";
            navigation.style.width = "100%";
            main.style.animation = "SlideUp3 1.2s ease-in";
        }

        /*  API  */
        const api_key = localStorage.getItem('api_key');
        $.ajaxSetup
        ({
            headers: { 'Authorization': 'ApiKey ' + api_key}
        });
        
        $.get("http://localhost:8080/api/customers", [], [])
        .done(function( _data) 
        {
            console.log(_data);
            setData(_data)
        })
        .fail( function(xhr) 
        {
            alert(xhr.responseText);
        });
        
    }, []);

    return (
        <>
            <div className = "main" style = {{left: '50%', top: '53%', transform: 'translate(-50%, -50%)', 
                                        height: '90%', backgroundColor: 'transparent', animation:'SlideUp3 1.2s ease-in'}}>
                <div className = "search">
                    <form className = "search-area">
                        <input className ="search-area" type="search" placeholder="Search..." />
                    </form>    
                </div>
                <div className = "main-elements">
                    {data.map((_data, id)=>
                        {
                            return <Customer_details />

                        })
                    }
                    <Customer_details />
                    <Customer_details />
                </div>
                <div className = "center" id = "pag"></div>
            </div>

            <Page1 filter_display = "none"/>
            
        </>
    );
}

export default Customers;