import {useEffect, useState} from 'react';
import Page1 from '../components/Page1';
import $ from 'jquery';
import Order_details from '../components/semi-components/order-details';
import '../CSS/page1.css';

/* Must start with a Caps letter */
function Orders()
{
    const[inputs, setInputs] = useState({});

    const handleChange = (event) =>
    {
        const name = event.target.name;
        const value = event.target.value;
        setInputs(values => ({...values, [name]: value}))
    }
    const [data, setData] = useState([]);

    const SearchOrder = (event) =>
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
        
        $.get("http://localhost:8080/api/orders", [], [])
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
        <div className = "orders">
            <div className = "main" style = {{left: '50%', top: '53%', transform: 'translate(-50%, -50%)', 
                                        height: '90%', backgroundColor: 'transparent', animation:'SlideUp3 1.2s ease-in'}}>

                <div className = "search" onSubmit={(event) => SearchOrder(event)}>
                    <form className = "search-area" autoComplete='off'>
                    <input className ="search-area" type="search" placeholder="Search..." 
                        name = "search" value = {inputs.search || ""}  onChange = {handleChange}></input>
                    </form>    
                </div>
                <div className = "main-elements">
                    <div className = "pan-main">
                        {data.map((_data, id)=>
                            {
                                return <Order_details />

                            })
                        }
                        <Order_details />
                        <Order_details />
                    </div>
                </div>
                <div className = "center" id = "pag"></div>
            </div>

            <Page1 filter_display = "none"/>
            
        </div>
    );
}

export default Orders;