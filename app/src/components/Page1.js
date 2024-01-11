import {useEffect, useState} from 'react';
import { flushSync } from 'react-dom';
import { createRoot } from 'react-dom/client';

import Detailed_Images from './semi-components/Product/detailed_images';
import Detailed_Images2 from './semi-components/Product/detailed_images2';
import Detailed_product from './semi-components/Product/detailed_product';
import Product_Variants from './semi-components/Product/product_variants';
import Detailed_Price from '../components/semi-components/Product/detailed_prices';
import Detailed_Quantities from '../components/semi-components/Product/detailed_quantities';
import Detailed_Options from '../components/semi-components/Product/detailed_options';

import $ from 'jquery';
import Pan_details from './semi-components/pan-detail';
import '../CSS/page1.css';


function Page1(props)
{
    const[inputs, setInputs] = useState({});

    const handleChange = (event) => 
    {
        const name = event.target.name;
        const value = event.target.value;
        setInputs(values => ({...values, [name]: value}))
      }

    const Filter = (event) =>
    {
        event.preventDefault();

        let category = document.querySelector(".category");
        let type = document.querySelector(".type");
        let vendor = document.querySelector(".vendor");

        if(inputs.category == undefined) {inputs.category = "";}
        if(inputs.vendor == undefined){inputs.vendor = ""; }
        if(inputs.type == undefined){inputs.type = ""; }

        category.innerHTML = inputs.category;
        type.innerHTML = inputs.type;
        vendor.innerHTML = inputs.vendor;

        let filter_input = document.querySelectorAll(".filter-selection-main");
        let navbar = document.querySelector(".navbar"); let main = document.querySelector(".main"); 
        let filter = document.querySelector(".filter");
        for(let i = 0; i < filter_input.length; i++) { filter_input[i].style.display = "none"; }
        let filter_button = document.getElementById("_filter"); let C_filter = document.getElementById("clear_filter");
        filter_button.disabled = false; C_filter.disabled = false;
        filter_button.style.cursor = "pointer"; C_filter.style.cursor = "pointer"; navbar.style.display = "block";
        main.style.display = "block"; filter.style.display = "block";
    }


    useEffect(()=> 
    {
        /* Ensure the model is shown */
        let navbar = document.getElementById("navbar");
        navbar.style.display = "block";
        
        /* animation for the search bar */
        let search = document.querySelector(".search-area");
        setTimeout(() =>
        { search.style.opacity = "1"; search.style.animation = "appear 0.8s ease-in"; }, 1000);

        /* animation for the pan elements */
        let pan = document.querySelectorAll(".pan");
        let pag = document.getElementById("pag");
        setTimeout(() =>
        {
            for(let i = 0; i < pan.length; i ++)
            {
                pan[i].style.display = "block";
                pan[i].style.animation = "appear 0.8s ease-in";
            }
            pag.style.display = "block";
            pag.style.animation = "appear 1s ease-in";
        }, 1000);


        /* filter element animation */
        let filters = document.querySelector(".filter").children;
        setTimeout(() =>
        {
            for(let i = 0; i < filters.length; i ++)
            {
                filters[i].style.display = "block";
                filters[i].style.animation = "appear 0.8s ease-in";
            }
        }, 1000);

        /* filter image script to show when clicked on */
        let filter_button = document.getElementById("_filter");
        let filter = document.querySelectorAll(".filter-elements");
        let filter_img = document.querySelectorAll(".filter-img");
        let C_filter = document.getElementById("clear_filter");
        let filter_input = document.querySelectorAll(".filter-selection-main");
        let close = document.querySelectorAll(".close-filter");


        filter_button.disabled = true;
        C_filter.disabled = true;

        for(let i = 0; i < filter.length; i++)
        {
            /* Filter Onclick */
            filter[i].addEventListener("click", () =>
            {
                filter_img[i].style.display = "block";
                filter[i].style.backgroundColor = "rgba(64, 165, 24, 0.7)";
                filter_input[i].style.display = "block";
            });

            /* Filter Close button onclick */
            close[i].addEventListener("click", () =>
            {
                filter_img[i].style.display = "none";
                filter[i].style.backgroundColor = "rgba(61, 61, 61, 0.7)";
                filter_input[i].style.display = "none";
            });

            /* Clear Filter */
            C_filter.addEventListener("click", () =>
            {
                filter_img[i].style.display = "none";
                filter[i].style.backgroundColor = "rgba(61, 61, 61, 0.7)";
            });
            
        }

        /* When the user clicks on the pan elements show info about that specified pan element */
        function DetailedView()
        {
            let products = document.querySelector(".products");
            let pan = document.querySelectorAll(".pan");
            for(let i = 0; i < pan.length; i++)
            {
                pan[i].addEventListener("click", () =>
                {
                    let id = pan[i].querySelector(".p-d-id").innerHTML;
                    /*  API  */
                    const api_key = localStorage.getItem('api_key');
                    $.ajaxSetup({ headers: { 'Authorization': 'ApiKey ' + api_key} });
                    $.get("http://localhost:8080/api/products/" + id, [], [], 'json')
                    .done(function(_data) 
                    {   
                        console.log(_data);
                        if(document.querySelector(".details") != null)
                        //div already exists, remove it, and create another
                        {

                            document.querySelector(".details").remove();
                            let details = document.createElement('div');
                            details.className = "details";
                            products.appendChild(details);

                            let rot = createRoot(details);
                            rot.render( <Detailed_product key={`${_data.title}_${i}`} Product_Title = {_data.title} Product_Category={_data.category} Product_Code={_data.product_code}
                                Product_Type={_data.product_type} Product_Vendor={_data.vendor} Product_ID={_data.id} Product_Activity={_data.active}
                                Product_Options={_data.options.map((el, i) => <Detailed_Options key={`${el.title}_${i}`} Option_Value={el.value} Option_Name = {el.position} />)}
                            />)
                            /* For some reason it wont pick up the element unless it throw it here */
                            setTimeout(() =>
                            {
                                details.querySelector(".description").innerHTML = _data.body_html;

                                let _div = details.querySelectorAll(".auto-slideshow-container");
                                for(let i = 0; i < _div.length; i++)
                                {
                                    let _root = createRoot(_div[i]);
                                    if(i == 0) { _root.render( _data.product_images.map((el, i) => <Detailed_Images key={`${el.title}_${i}`} Image1 = {el.src}/> )) }
                                    else { _root.render( _data.product_images.map((el, i) => <Detailed_Images2 key={`${el.title}_${i}`} Image1 = {el.src}/> )) }
                                }
                                let new_div = details.querySelector(".variants"); 
                                let rt = createRoot(new_div);
                                rt.render( _data.variants.map((el, i) => <Product_Variants key={`${el.title}_${i}`} Variant_Title = {el.id}
                                Variant_Barcode={el.barcode} Variant_SKU={el.sku} Variant_UpdateDate={el.updated_at} 
                                Option1={el.option1} Option2={el.option2} Option3={el.option3} 
                                Price={el.variant_price_tiers.map((el, i) => <Detailed_Price key={`${el.title}_${i}`} Price_Name={el.name} Price_Value={el.value}  />)}
                                Quantities={el.variant_quantities.map((el, i) => <Detailed_Quantities quantity_value = {el.value}/>)}
                                />))
                            }, 10);  
                        }
                        else 
                        //create new div
                        {
                            let details = document.createElement('details');
                            products.appendChild(details);
                            let rot = createRoot(details);
                            rot.render( <Detailed_product key={`${_data.title}_${i}`} Product_Title = {_data.title} Product_Category={_data.category} Product_Code={_data.product_code}
                                Product_Type={_data.product_type} Product_Vendor={_data.vendor} Product_ID={_data.id} Product_Activity={_data.active}
                                Product_Options={_data.options.map((el, i) => <Detailed_Options key={`${el.title}_${i}`} Option_Value={el.value} Option_Name = {el.position} />)}
                            />)
                            /* For some reason it wont pick up the element unless it throw it here */
                            setTimeout(() =>
                            {
                                details.querySelector(".description").innerHTML = _data.body_html;
                                let _div = details.querySelectorAll(".auto-slideshow-container");
                                for(let i = 0; i < _div.length; i++)
                                {
                                    let _root = createRoot(_div[i]);
                                    if(i == 0)
                                    {
                                        _root.render( _data.product_images.map((el, i) => <Detailed_Images key={`${el.title}_${i}`} Image1 = {el.src}/> ))
                                    }
                                    else 
                                    {
                                        _root.render( _data.product_images.map((el, i) => <Detailed_Images2 key={`${el.title}_${i}`} Image1 = {el.src}/> ))
                                    }
                                }
                                let new_div = details.querySelector(".variants"); 
                                let rt = createRoot(new_div);
                                rt.render( _data.variants.map((el, i) => <Product_Variants key={`${el.title}_${i}`} Variant_Title = {el.id}
                                Variant_Barcode={el.barcode} Variant_SKU={el.sku} Variant_UpdateDate={el.updated_at} 
                                Option1={el.option1} Option2={el.option2} Option3={el.option3} 
                                Price={el.variant_price_tiers.map((el, i) => <Detailed_Price key={`${el.title}_${i}`} Price_Name={el.name} Price_Value={el.value}  />)}
                                Quantities={el.variant_quantities.map((el, i) => <Detailed_Quantities quantity_value = {el.value}/>)}
                                /> ))
                            }, 0);
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

        /* Filter */
        filter_button.addEventListener("click", () =>
        {
            let category = document.querySelector(".category").innerHTML;
            let type = document.querySelector(".type").innerHTML;
            let vendor = document.querySelector(".vendor").innerHTML;
            let next = document.getElementById("next");

            $.get("http://localhost:8080/api/products/filter?type=" +type + "&" + "vendor="+ vendor +"&category="+category,[], [], 'json')
            .done(function( _data) 
            {
                console.log(_data);
                
                if(_data.length < 10)
                {
                    next.disabled = true;
                    next.style.cursor = "not-allowed";
                }
                if(_data == "")
                {
                    document.querySelector(".pan-main").remove();
                    document.querySelector(".empty-message").style.display = "block";
                }
                else 
                {
                    if(document.querySelector(".pan-main") != null)
                    {
                        document.querySelector(".pan-main").remove();
                        let pan_main = document.createElement('div');
                        let main_elements = document.querySelector(".main-elements");
                        pan_main.className = "pan-main";
                        main_elements.appendChild(pan_main);

                        let root = createRoot(pan_main);
                        flushSync(() => 
                        {
                            root.render(_data.map((el, i) => <Pan_details key={`${el.title}_${i}`} Product_Title={el.title} Product_ID={el.id} Product_Activity={el.active}
                            Product_Type={el.product_type} Product_Code={el.product_code} Product_Category={el.category} Product_Vendor={el.vendor}
                            Product_Image={el.product_images.map((el, i) => el.src)}
                            /> )) 
                        });
                        setTimeout(() => { DetailedView();}, 200);
                    }
                    else 
                    {
                        let pan_main = document.createElement('div');
                        let main_elements = document.querySelector(".main-elements");
                        pan_main.className = "pan-main";
                        main_elements.appendChild(pan_main);

                        let root = createRoot(pan_main);
                        flushSync(() => 
                        {
                            root.render(_data.map((el, i) => <Pan_details key={`${el.title}_${i}`} Product_Title={el.title} Product_ID={el.id} Product_Activity={el.active}
                            Product_Type={el.product_type} Product_Code={el.product_code} Product_Category={el.category} Product_Vendor={el.vendor}
                            Product_Image={el.product_images.map((el, i) => el.src)}
                            /> )) 
                        });
                        setTimeout(() => { DetailedView();}, 200);
                    }
                } 
            })
            .fail( function(xhr) { alert(xhr.responseText); });
        });

        /* Script to automatically format the number of elements on each page */
        const content = document.querySelector('.center'); 
        const paginationContainer = document.createElement('div');
        const paginationDiv = document.body.appendChild(paginationContainer);
        paginationContainer.classList.add('pagination');
        content.appendChild(paginationContainer);

        /* paginiation function for filtered elements */
        function Filter_Pagintation(index)
        {
            if(index == 1)
            {
                let category = document.querySelector(".category").innerHTML;
                let type = document.querySelector(".type").innerHTML;
                let vendor = document.querySelector(".vendor").innerHTML;
                let ahead = index + 1;
                $.get("http://localhost:8080/api/products/filter?type=" + type + "&vendor=" + vendor + "&category=" + category + "&page=" + ahead, [], [])
                .done(function( _data) 
                {
                    if(_data == "") { let next = document.getElementById("next"); next.style.cursor = "not-allowed"; next.disabled = true; } 
                })
                .fail( function(xhr) { alert(xhr.responseText); });
            }

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
                    let category = document.querySelector(".category").innerHTML;
                    let type = document.querySelector(".type").innerHTML;
                    let vendor = document.querySelector(".vendor").innerHTML;

                    index = index + 1;
                    /*  API  */
                    const api_key = localStorage.getItem('api_key');
                    $.ajaxSetup({ headers: { 'Authorization': 'ApiKey ' + api_key} });
                    $.get("http://localhost:8080/api/products/filter?type=" + type + "&vendor=" + vendor + "&category=" + category + "&page=" + index, [], [])
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
                            root.render(_data.map((el, i) => <Pan_details key={`${el.title}_${i}`} Product_Title={el.title} Product_ID={el.id} Product_Activity={el.active}
                            Product_Type={el.product_type} Product_Code={el.product_code} Product_Category={el.category} Product_Vendor={el.vendor}
                            Product_Image={el.product_images.map((el, i) => el.src)}
                            /> )) 
                        });
                        
                    })
                    .fail( function(xhr) { alert(xhr.responseText); });

                    let ahead = index + 1;
                    /*  API  */
                    $.get("http://localhost:8080/api/products/filter?type=" + type + "&vendor=" + vendor + "&category=" + category + "&page=" + ahead, [], [])
                    .done(function( _data) 
                    {
                        console.log(_data);
                        if(_data == "") { let next = document.getElementById("next"); next.style.cursor = "not-allowed"; next.disabled = true; }  
                    })
                    .fail( function(xhr) { alert(xhr.responseText); });

                    Filter_Pagintation(index);
                    setTimeout(() => { DetailedView();}, 200);
                });

    
                prevPage.addEventListener("click", () =>
                {
                    let category = document.querySelector(".category").innerHTML;
                    let type = document.querySelector(".type").innerHTML;
                    let vendor = document.querySelector(".vendor").innerHTML;
                    index = index - 1;
                    const api_key = localStorage.getItem('api_key');
                    $.ajaxSetup({ headers: { 'Authorization': 'ApiKey ' + api_key} });
                    $.get("http://localhost:8080/api/products/filter?type=" + type + "&vendor=" + vendor + "&category=" + category + "&page=" + index, [], [])
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
                            root.render(_data.map((el, i) => <Pan_details key={`${el.title}_${i}`} Product_Title={el.title} Product_ID={el.id} Product_Activity={el.active}
                            Product_Type={el.product_type} Product_Code={el.product_code} Product_Category={el.category} Product_Vendor={el.vendor}
                            Product_Image={el.product_images.map((el, i) => el.src)}
                            /> )) 
                        });

                    
                    })
                    .fail( function(xhr) { alert(xhr.responseText); });
                    Filter_Pagintation(index);
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
                    let category = document.querySelector(".category").innerHTML;
                    let type = document.querySelector(".type").innerHTML;
                    let vendor = document.querySelector(".vendor").innerHTML;

                    index = index + 1;

                    const api_key = localStorage.getItem('api_key');
                    $.ajaxSetup({ headers: { 'Authorization': 'ApiKey ' + api_key} });
                    $.get("http://localhost:8080/api/products/filter?type=" + type + "&vendor=" + vendor + "&category=" + category + "&page=" + index, [], [])
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
                            root.render(_data.map((el, i) => <Pan_details key={`${el.title}_${i}`} Product_Title={el.title} Product_ID={el.id} Product_Activity={el.active}
                            Product_Type={el.product_type} Product_Code={el.product_code} Product_Category={el.category} Product_Vendor={el.vendor}
                            Product_Image={el.product_images.map((el, i) => el.src)}
                            /> )) 
                        });

                    })
                    .fail( function(xhr) { alert(xhr.responseText); });
                    Filter_Pagintation(index);
                    setTimeout(() => { DetailedView();}, 200);
                });

                prevPage.addEventListener("click", () =>
                {
                    let category = document.querySelector(".category").innerHTML;
                    let type = document.querySelector(".type").innerHTML;
                    let vendor = document.querySelector(".vendor").innerHTML;
                    index = index - 1;

                    const api_key = localStorage.getItem('api_key');
                    $.ajaxSetup({ headers: { 'Authorization': 'ApiKey ' + api_key} });
                    $.get("http://localhost:8080/api/products/filter?type=" + type + "&vendor=" + vendor + "&category=" + category + "&page=" + index, [], [])
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
                            root.render(_data.map((el, i) => <Pan_details key={`${el.title}_${i}`} Product_Title={el.title} Product_ID={el.id} Product_Activity={el.active}
                            Product_Type={el.product_type} Product_Code={el.product_code} Product_Category={el.category} Product_Vendor={el.vendor}
                            Product_Image={el.product_images.map((el, i) => el.src)}
                            /> ))  
                        });
                        
                    })
                    .fail( function(xhr) { alert(xhr.responseText); });
                    Filter_Pagintation(index);
                    setTimeout(() => { DetailedView();}, 200);
                });
            } 
        }

        filter_button.addEventListener("click", () =>
        {
            Filter_Pagintation(1);
        });

        /* Fetch-Products Button */
        let fetch = document.getElementById("fetch");
        fetch.addEventListener("click", () =>
        {
            fetch.disabled = true;
            fetch.style.cursor = "not-allowed";

            let message = document.getElementById("message");
            message.style.display = "block";

            const api_key = localStorage.getItem('api_key');
            $.ajaxSetup({ headers: { 'Authorization': 'ApiKey ' + api_key} });
            $.get("http://localhost:8080/api/shopify/fetch", [], [])
            .done(function( _data) 
            {
                console.log(_data);
                fetch.disabled = false;
                fetch.style.cursor = "pointer";

                message.innerHTML = "success";
                message.style.background = "#1a5e12";
                setTimeout(() =>
                {
                    message.innerHTML = "";
                    message.style.backgroundColor = "transparent";
                    message.style.display = "none";
                }, 2000);
            })
            .fail( function(xhr) { alert(xhr.responseText); }); 
            
        });


        /* Push(SYNC)-Products Button */
        let push = document.getElementById("push");
        push.addEventListener("click", () =>
        {
            push.disabled = true;
            push.style.cursor = "not-allowed";

            let message = document.getElementById("message");
            message.style.display = "block";

            const api_key = localStorage.getItem('api_key');
            $.ajaxSetup({ headers: { 'Authorization': 'ApiKey ' + api_key} });
            $.get("http://localhost:8080/api//shopify/sync", [], [])
            .done(function( _data) 
            {
                console.log(_data);
                push.disabled = false;
                push.style.cursor = "pointer";

                message.innerHTML = _data.message;
                message.style.background = "#1a5e12";
                setTimeout(() =>
                {
                    message.innerHTML = "";
                    message.style.backgroundColor = "transparent";
                    message.style.display = "none";
                }, 2000);
            })
            .fail( function(xhr) { alert(xhr.responseText); }); 
        });
        
        
        


    }, []);

    return (
        <>
            <div className = "filter" style = {{display: props.filter_display}}>
                <div className = "filter-title"><b>Available Filters:</b></div>
                <div className = "filter-elements">
                    <div className = "filter-elements-text">Filter By Type</div>
                    <div className = "filter-img"/>
                    <div className = "type"></div>
                </div>
                <div className = "filter-elements">
                    <div className = "filter-elements-text">Filter By Type</div>
                    <div className = "filter-img"/>
                    <div className = "vendor"></div>
                </div>
                <div className = "filter-elements">
                    <div className = "filter-elements-text">Filter By Type</div>
                    <div className = "filter-img"/>
                    <div className = "category"></div>
                </div>
                <br />
                <div className = "vendor"></div>
                <div className = "type"></div>
                <div className = "category"></div>

                <div id = "button-hold">
                    <button id = "clear_filter"className = "filter-button">Clear Filter</button>
                    <button id = "_filter"className = "filter-button">Filter Results</button>
                </div>
                
                <br/><br/>
                <div className = "filter-title"><b>Product Sync</b></div>
                <div className = "fetch-product" id = "fetch"><div className = "filter-elements-text">Fetch Products</div></div>
                <div className = "push-product" id = "push"><div className = "filter-elements-text">Push Products</div></div>

                
            </div>
            <div className = "filter-selection-main">
                <div className = "filter-input">
                    <div className = 'close-filter'>&times;</div>
                    <div className = "filter-selection-title">Filter Type</div>
                    <form method = 'post' onSubmit={(event) => Filter(event)} autoComplete='off'>
                        <span><input type = 'type' placeholder = "Enter Type" name = "type" value = {inputs.type || ""} onChange = {handleChange} required></input></span>
                        <br/><br/><br/>
                        <button className = 'button' type = 'submit'>Confirm</button>
                    </form>

                </div>
            </div>
            <div className = "filter-selection-main">
                <div className = "filter-input">
                    <div className = 'close-filter'>&times;</div>
                    <div className = "filter-selection-title">Filter Vendor</div>
                    <form method = 'post' onSubmit={(event) => Filter(event)} autoComplete='off'>
                        <span><input type = 'vendor' placeholder = "Enter Vendor" name = "vendor" value = {inputs.vendor || ""} onChange = {handleChange} required></input></span>
                        <br/><br/><br/>
                        <button className = 'button' type = 'submit'>Confirm</button>
                    </form>
                </div>
            </div>
            <div className = "filter-selection-main">
                <div className = "filter-input">
                    <div className = 'close-filter'>&times;</div>
                    <div className = "filter-selection-title">Filter Category</div>
                    <form method = 'post' onSubmit={(event) => Filter(event)} autoComplete='off'>
                        <span><input type = 'type' placeholder = "Enter Category" name = "category" value = {inputs.category || ""} onChange = {handleChange} required></input></span>
                        <br/><br/><br/>
                        <button className = 'button' type = 'submit'>Confirm</button>
                    </form>
                </div>
            </div>
            <div className = 'info-message' id = 'message' />
        </>
    );
}
Page1.defaultProps = 
{
    filter_display: 'block', 
    main_display: 'block',
    main_bgc: '',
    main_top: '13%',
    main_left: '51%',
    transform: 'translate(-30%, -6%)',
    width: '70%',
    height: '96%', 
    animation: 'SlideUp2 0.8s ease-in'
};

export default Page1;