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

        /* When the user clicks on the pan elements show info about that specified pan element */
        let pan = document.querySelectorAll(".pan");
        for(let i = 0; i < pan.length; i++)
        {
            pan[i].addEventListener("click", () =>
            {
                
                //var img = pan[i].querySelector(".pan-img").innerHTML; 
                document.getElementById("img").style.backgroundImage = pan[i].querySelector(".pan-img").style.backgroundImage;
                document.getElementById("te").innerHTML = pan[i].querySelector(".p-d-title").innerHTML;
                document.getElementById("co").innerHTML = pan[i].querySelector(".p-d-code").innerHTML;
                document.getElementById("op").innerHTML = pan[i].querySelector(".p-d-options").innerHTML; 
                document.getElementById("ca").innerHTML = pan[i].querySelector(".p-d-category").innerHTML;
                document.getElementById("ty").innerHTML = pan[i].querySelector(".p-d-type").innerHTML; 
                document.getElementById("ve").innerHTML = pan[i].querySelector(".p-d-vendor").innerHTML;

                /* Get the filter & main elements */
                let filter = document.querySelector(".filter");
                let main = document.querySelector(".main");
                let navbar = document.getElementById("navbar");
                let details = document.querySelector(".details");
                let close = document.querySelector(".close-button");

                filter.style.animation = "Fadeout 0.5s ease-out";
                main.style.animation = "Fadeout 0.5s ease-out";
                navbar.style.animation = "Fadeout 0.5s ease-out";

                setTimeout(() => 
                {
                    filter.style.display = "none";
                    main.style.display = "none";
                    navbar.style.display = "none";

                    details.style.animation = "FadeIn ease-in 0.5s";
                    details.style.display = "block";
                    close.style.display = "block";
                }, 500);
            });

            /* When the user clicks on the return button */
            let close = document.querySelector(".close-button");
            let filter = document.querySelector(".filter");
            let main = document.querySelector(".main");
            let navbar = document.getElementById("navbar");
            let details = document.querySelector(".details");
            close.addEventListener("click", ()=> 
            {
                close.style.display = "none";
                details.style.animation = "Fadeout 0.5s ease-out";
                main.style.animation = "FadeIn ease-in 0.5s";
                filter.style.animation = "FadeIn ease-in 0.5s";
                navbar.style.animation = "FadeIn ease-in 0.5s";

                details.style.display = "none";
                navbar.style.display = "block";
                main.style.display = "block";
            });
        }
        
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
            <div className = "details">
                <div className = 'close-button'>&times;</div>
                <div id = "img" className = "details-image">
                    <div id = "te" className = "details-details details-title"></div>
                    <div id = "co" className = "details-details details-code"></div>
                    <div id = "op" className = "details-details details-options"></div>
                    <div id = "ca" className = "details-details details-category"></div>
                    <div id = "ty" className = "details-details details-type"></div>
                    <div id = "ve" className = "details-details details-vendor"></div>
                </div>
                
            </div>
            
        </div>
    );
}

export default Customers;