import { createRoot } from 'react-dom/client';
import { flushSync } from 'react-dom';
import {useEffect, useState} from 'react';
import ReactDOM from 'react-dom/client';
import $ from 'jquery';

import Page1 from '../components/Page1';
import Pan_details from '../components/semi-components/pan-detail';
import Detailed_product from '../components/semi-components/Product/detailed_product';
import Product_Variants from '../components/semi-components/Product/product_variants';
import Detailed_Images from '../components/semi-components/Product/detailed_images';
import product from '../media/products.png';

import '../CSS/page1.css';

function Products()
{
    const[inputs, setInputs] = useState({});

    const handleChange = (event) =>
    {
        const name = event.target.name;
        const value = event.target.value;
        setInputs(values => ({...values, [name]: value}))
    }

    const [data, setData] = useState([]);
    //const [data2, setData2] = useState([]);

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
            navigation.style.animation = "MoveLeft 0.8s ease";
        }

        /*  API  */
        const api_key = localStorage.getItem('api_key');
        $.ajaxSetup({ headers: { 'Authorization': 'ApiKey ' + api_key} });
        $.get("http://localhost:8080/api/products?page=1", [], [])
        .done(function( _data) 
        {
            console.log(_data);
            setData(_data);
        })
        .fail( function(xhr) 
        {
            alert(xhr.responseText);
        });


        /* When the user clicks on the pan elements show info about that specified pan element */
        function DetailedView()
        {
            let pan = document.querySelectorAll(".pan");
            console.log("pan");
            
            for(let i = 0; i < pan.length; i++)
            {
                pan[i].addEventListener("click", () =>
                {
                    console.log([i] + " was clicked");
                    let id = pan[i].querySelector(".p-d-id").innerHTML;

                    /*  API  */
                    const api_key = localStorage.getItem('api_key');
                    $.ajaxSetup({ headers: { 'Authorization': 'ApiKey ' + api_key} });
                    $.get("http://localhost:8080/api/products/" + id, [], [], 'json')
                    .done(function(_data) 
                    {
                        console.log(_data);
                        //setData2(_data);


                        
                        let details = document.querySelector(".details");
                        let rot = createRoot(details);

                        rot.render( <Detailed_product Product_Title = {_data.title} />)
                        

                        /* For some reason it wont pick up the element unless it throw it here */
                        setTimeout(() =>
                        {
                            let _div = details.querySelector(".auto-slideshow-container");
                            let _root = createRoot(_div);

                            _root.render( _data.product_images.map((el, i) =>
                                <Detailed_Images key={`${el.title}_${i}`} Image1 = {el.src}/>
                            ))

                            let new_div = details.querySelector(".variants"); 
                            let rt = createRoot(new_div);
    
                            rt.render( _data.variants.map((el, i) =>
                                <Product_Variants key={`${el.title}_${i}`} Variant_Title = {el.id}/>
                            ))

                            /*
                                root.render(_data.map((el, i) => 
                                <Pan_details key={`${el.title}_${i}`} Product_Title={el.title} Product_ID={el.id}/>
                                ))
                            */
                        }, 50);
                         

                           
                        
                        /*
                           
                        */

                        /*
                        {data.map((el, i) => 
                            <Pan_details key={`${el.title}_${i}`} Product_Title={el.title} Product_ID={el.id}
                            Product_Code={el.product_code}
                            />
                        )}
                        */

                    })
                    .fail( function(xhr) 
                    {
                        alert(xhr.responseText);
                    });


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
            
        }

        /* Script to automatically format the number of elements on each page */
        const content = document.querySelector('.center'); 
        const paginationContainer = document.createElement('div');
        const paginationDiv = document.body.appendChild(paginationContainer);
        paginationContainer.classList.add('pagination');
        content.appendChild(paginationContainer);

        let div = document.getElementById("pan-main");
        let root = createRoot(div);

        function Pagintation(index)
        {
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
                if(index == 1)
                {
                    prevPage.disabled = true;
                    prevPage.style.cursor = "not-allowed";
                }
                else if(index > 1)
                {
                    prevPage.style.cursor = "pointer";
                    prevPage.disabled = false;
                    nextPage.disabled = false;
                }
                else if(index <= 1)
                {
                    prevPage.disabled = true;
                    prevPage.style.cursor = "not-allowed";
                }

                nextPage.addEventListener("click", () =>
                {
                    index = index + 1;
                    /* Fetches the data from page, based on the page / index value */
                    const page = "http://localhost:8080/api/products?page=" + index;
                    /*  API  */
                    const api_key = localStorage.getItem('api_key');
                    $.ajaxSetup({ headers: { 'Authorization': 'ApiKey ' + api_key} });
                    $.get(page, [], [])
                    .done(function( _data) 
                    {
                        console.log(_data);

                        flushSync(() => 
                        {
                            root.render(_data.map((el, i) => 
                                <Pan_details key={`${el.title}_${i}`} Product_Title={el.title} Product_ID={el.id}/>
                            ))
                        });
                        
                    })
                    .fail( function(xhr) 
                    {
                        alert(xhr.responseText);
                    });

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
                    .fail( function(xhr) 
                    {
                        alert(xhr.responseText);
                    });

                    Pagintation(index);
                    setTimeout(() => { DetailedView(); }, 500);
                });

                prevPage.addEventListener("click", () =>
                {
                    index = index - 1;
                    /* Fetches the data from page, based on the page / index value */
                    const page = "http://localhost:8080/api/products?page=" + index;

                    /*  API  */
                    const api_key = localStorage.getItem('api_key');
                    $.ajaxSetup({ headers: { 'Authorization': 'ApiKey ' + api_key} });
                    $.get(page, [], [])
                    .done(function( _data) 
                    {
                        console.log(_data);
                        flushSync(() => 
                        {
                            root.render(_data.map((el, i) => 
                                <Pan_details key={`${el.title}_${i}`} Product_Title={el.title} Product_ID={el.id}/>
                            ))
                        });
                    })
                    .fail( function(xhr) 
                    {
                        alert(xhr.responseText);
                    });

                    Pagintation(index--);
                    setTimeout(() => { DetailedView(); }, 500);
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

                if(index == 1)
                {
                    prevPage.disabled = true;
                    prevPage.style.cursor = "not-allowed";
                }
                else if(index > 1)
                {
                    prevPage.style.cursor = "pointer";
                    prevPage.disabled = false;
                    nextPage.disabled = false;
                }
                else if(index <= 1)
                {
                    prevPage.disabled = true;
                    prevPage.style.cursor = "not-allowed";
                }
                nextPage.addEventListener("click", () =>
                {
                    index = index + 1;
                    /* Fetches the data from page, based on the page / index value */
                    const page = "http://localhost:8080/api/products?page=" + index;
                    /*  API  */
                    const api_key = localStorage.getItem('api_key');
                    $.ajaxSetup({ headers: { 'Authorization': 'ApiKey ' + api_key} });
                    $.get(page, [], [])
                    .done(function( _data) 
                    {
                        console.log(_data);
                        /* Check if pan elements exist and remove + update if it does*/
                        let pan = document.querySelectorAll(".pan");
                        pan.forEach(pan => { pan.remove(); });
                        flushSync(() => 
                        {
                            root.render(_data.map((el, i) => 
                                <Pan_details key={`${el.title}_${i}`} Product_Title={el.title} Product_ID={el.id}/>
                            ))
                        });
                    })
                    .fail( function(xhr) 
                    {
                        alert(xhr.responseText);
                    });

                    Pagintation(index++);
                    setTimeout(() => { DetailedView(); }, 500);
                });

                prevPage.addEventListener("click", () =>
                {
                    index = index - 1;
                    /* Fetches the data from page, based on the page / index value */
                    const page = "http://localhost:8080/api/products?page=" + index;

                    /*  API  */
                    const api_key = localStorage.getItem('api_key');
                    $.ajaxSetup({ headers: { 'Authorization': 'ApiKey ' + api_key} });
                    $.get(page, [], [])
                    .done(function( _data) 
                    {
                        console.log(_data);
                        /* Check if pan elements exist and remove + update if it does*/
                        let pan = document.querySelectorAll(".pan");
                        pan.forEach(pan => { pan.remove(); });
                        flushSync(() => 
                        {
                            root.render(_data.map((el, i) => 
                                <Pan_details key={`${el.title}_${i}`} Product_Title={el.title} Product_ID={el.id}/>
                            ))
                        });
                    })
                    .fail( function(xhr) 
                    {
                        alert(xhr.responseText);
                    });

                    Pagintation(index);
                    setTimeout(() => { DetailedView(); }, 500);
                });
            } 
        }
        Pagintation(1);
        setTimeout(() => { DetailedView(); }, 500);
        

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
                    <div className = "pan-main" id = "pan-main">
                        {data.map((el, i) => 
                            <Pan_details key={`${el.title}_${i}`} Product_Title={el.title} Product_ID={el.id}
                            Product_Code={el.product_code}
                            />
                        )}
                    </div>
                </div>
                <div className = "center" id = "pag"></div>
            </div>

            <Page1 image = {product} title = "Products"/>
            <div className = "details">
                <div className = 'close-button'>&times;</div>
                
                
                       

            </div>

        </div>
    );
}

export default Products;

/*

{data2.map((el, i) => 
                    {
                        <Detailed_product Product_Title = {el.title} key={`${el.id}_${i}`} inputRef = 
                        {
                            <Product_Variants/>
                        } 
                    /> 
                })}


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

*/