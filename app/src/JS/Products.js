import {useEffect, useState} from 'react';
import $ from 'jquery';
import Page1 from '../components/Page1';
import Pan_details from '../components/semi-components/pan-detail';
import '../CSS/page1.css';
import product from '../media/products.png';

function Products(props)
{
    const[inputs, setInputs] = useState({});

    const handleChange = (event) =>
    {
        const name = event.target.name;
        const value = event.target.value;
        setInputs(values => ({...values, [name]: value}))
    }

    const [data, setData] = useState([]);

    const SearchProduct = (event) =>
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
        /* Ensures the page elements are set correctly */
        let navigation = document.getElementById("navbar");
        window.onload = function(event)
        {
            navigation.style.left = "30%";
            navigation.style.position = "absolute";
            navigation.style.width = "70%";
            navigation.style.animation = "MoveLeft 1.2s ease";
        }

        /*  API  */
        const api_key = localStorage.getItem('api_key');
        $.ajaxSetup
        ({
            headers: { 'Authorization': 'ApiKey ' + api_key}
        });
        $.get("http://localhost:8080/api/products", [], [])
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
                document.getElementById("pr").innerHTML = pan[i].querySelector(".pan-price").innerHTML;

                /* Get the filter & main elements */
                let filter = document.querySelector(".filter");
                let main = document.querySelector(".main");
                let navbar = document.getElementById("navbar");
                let details = document.querySelector(".details");
                let close = document.querySelector(".close-button");

                filter.style.animation = "Fadeout 0.5s ease-out";
                main.style.animation = "Fadeout 0.5s ease-out";
                navbar.style.animation = "Fadeout 0.5s ease-out";

                filter.style.display = "none";
                main.style.display = "none";
                navbar.style.display = "none";
                details.style.animation = "FadeIn ease-in 0.5s";
                details.style.display = "block";
                close.style.display = "block";

               
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
                setTimeout(() => 
                {
                    details.style.display = "none";
                    navbar.style.display = "block";
                    main.style.display = "block";
                    filter.style.display = "block";
                }, 500);
            });
        }
        
        
    }, []);

    return (
        <div className = "products">
            <div className = "main">
                <div className = "search">
                    <form className = "search-area" autoComplete='off' onSubmit={(event) => SearchProduct(event)}>
                        <input className ="search-area" type="search" placeholder="Search..." 
                        name = "search" value = {inputs.search || ""}  onChange = {handleChange}></input>
                    </form>    
                </div>
                <div className = "main-elements">
                    <div className = "pan-main">
                        {data.map((_data, id)=>
                            {
                                return <Pan_details />

                            })
                        }
                        <Pan_details Product_Title = "5-star sword" Product_Code = "#w123d" Product_Options = "True-False" Product_Category = "Gacha"
                        Product_Type = "SSR" Product_Vendor = "HottaGames" Product_Price = "$15"/>

                        <Pan_details Product_Title = "5-star " Product_Code = "#rf34g" Product_Options = "white black" Product_Category = "pog"
                        Product_Type = "SW@" Product_Vendor = "PMdfg" Product_Price = "$155"/>

                        <Pan_details Product_Title = "5 sword" Product_Code = "#kn39c" Product_Options = "542/544" Product_Category = "Posxc"
                        Product_Type = "postman" Product_Vendor = "keyboard" Product_Price = "$147"/>

                    </div>
                </div>
                <div className = "center" id = "pag">
                    
                </div>
            </div>

            <Page1 image = {product} title = "Products"/>
            <div className = "details">
                <div className = 'close-button'>&times;</div>
                <div id = "img" className = "details-image">
                    <div id = "te" className = "details-details details-title"></div>
                    <div id = "co" className = "details-details details-code"></div>
                    <div id = "op" className = "details-details details-options"></div>
                    <div id = "ca" className = "details-details details-category"></div>
                    <div id = "ty" className = "details-details details-type"></div>
                    <div id = "ve" className = "details-details details-vendor"></div>
                    <div id = "pr" className = "details-details details-price"></div>
                </div>
                
            </div>

        </div>
    );
}

export default Products;