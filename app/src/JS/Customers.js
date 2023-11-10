import {useEffect, useState} from 'react';
import $ from 'jquery';
import Customer_details from '../components/semi-components/customer-details';
import Page1 from '../components/Page1';
import '../CSS/page1.css';

/* Must start with a Caps letter */
function Customers()
{
    const[inputs, setInputs] = useState({});

    const handleChange = (event) =>
    {
        const name = event.target.name;
        const value = event.target.value;
        setInputs(values => ({...values, [name]: value}))
    }
    const [data, setData] = useState([]);

    const SearchCustomer = (event) =>
    {
        event.preventDefault();
        console.log(inputs);

        /*
        $.post("http://localhost:8080/api/login", JSON.stringify(inputs),[], 'json')
        .done(function( _data) 
        {
            console.log(_data);
        })
        .fail( function(xhr) 
        {
            alert(xhr.responseText);
        });
        */
    }

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
        <div className = "customer">
            <div className = "main" style = {{left: '50%', top: '53%', transform: 'translate(-50%, -50%)', 
                                        height: '90%', backgroundColor: 'transparent', animation:'SlideUp3 1.2s ease-in'}}>
                <div className = "search">
                    <form className = "search-area" autoComplete = 'off' onSubmit={(event) => SearchCustomer(event)}>
                    <input className ="search-area" type="search" placeholder="Search..." 
                        name = "search" value = {inputs.search || ""}  onChange = {handleChange}></input>
                    </form>    
                </div>
                <div className = "main-elements">
                    <div className = "pan-main">
                        {data.map((_data, id)=>
                            {
                                return <Customer_details />

                            })
                        }
                        <Customer_details />
                        <Customer_details />
                    </div>
                </div>
                <div className = "center" id = "pag"></div>
            </div>

            <Page1 filter_display = "none"/>
            
        </div>
    );
}

export default Customers;