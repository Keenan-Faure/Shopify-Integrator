import {useEffect, useState} from 'react';
import { flushSync } from 'react-dom';
import { createRoot } from 'react-dom/client';
import Background from './Background';
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
        for(let i = 0; i < filter_input.length; i++)
        {
            filter_input[i].style.display = "none";
        }
    }


    useEffect(()=> 
    {
        /* Ensure the model is shown */
        let navbar = document.getElementById("navbar");
        navbar.style.display = "block";
        
        /* animation for the search bar */
        let search = document.querySelector(".search-area");
        setTimeout(() =>
        {
            search.style.opacity = "1";
            search.style.animation = "appear 0.8s ease-in";
        }, 1000);

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
        let filter = document.querySelectorAll(".filter-elements");
        let filter_main = document.querySelector(".filter");
        let filter_img = document.querySelectorAll(".filter-img");
        let C_filter = document.querySelector(".filter-button");
        let filter_input = document.querySelectorAll(".filter-selection-main");
        let close = document.querySelectorAll(".close-filter");
        let main = document.querySelector(".main");

        for(let i = 0; i < filter.length; i++)
        {
            /* Filter Onclick */
            filter[i].addEventListener("click", () =>
            {
                filter_img[i].style.display = "block";
                filter[i].style.backgroundColor = "rgba(64, 165, 24, 0.7)";
                filter_input[i].style.display = "block";
            });

            /* Clear Filter */
            C_filter.addEventListener("click", () =>
            {
                filter_img[i].style.display = "none";
                filter[i].style.backgroundColor = "rgba(61, 61, 61, 0.7)";
            });

            close[i].addEventListener("click", () =>
            {
                filter_img[i].style.display = "none";
                filter[i].style.backgroundColor = "rgba(61, 61, 61, 0.7)";
                filter_input[i].style.display = "none";
            });
        }

        /* Filter */
        let filter_button = document.getElementById("_filter");
        filter_button.addEventListener("click", () =>
        {
            let category = document.querySelector(".category").innerHTML;
            let type = document.querySelector(".type").innerHTML;
            let vendor = document.querySelector(".vendor").innerHTML;

            $.get("http://localhost:8080/api/products/filter?type=" +type + "&" + "vendor="+ vendor +"&category="+category,[], [], 'json')
            .done(function( _data) 
            {
                let div = document.getElementById("pan-main");
                let root = createRoot(div);
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
            
        });
        


    }, []);

    return (
        <>
            <Background />
            <div className = "filter" style = {{display: props.filter_display}}>
                <div className = "filter-title"><b>Available Filters:</b></div>

                <div className = "filter-elements">
                    Filter By Type
                    <div className = "filter-img"/>
                    <div className = "type"></div>
                </div>

                <div className = "filter-elements">
                    Filter By Vendor
                    <div className = "filter-img"/>
                    <div className = "vendor"></div>
                </div>

                <div className = "filter-elements">
                    Filter By Category
                    <div className = "filter-img"/>
                    <div className = "category"></div>
                </div>
                <br />

                <div className = "vendor"></div>
                <div className = "type"></div>
                <div className = "category"></div>

                <button className = "filter-button">Clear Filter</button>
                <button id = "_filter"className = "filter-button">Filter Results</button>
                
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

/*

*/