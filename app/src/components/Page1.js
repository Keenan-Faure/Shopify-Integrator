import {useEffect} from 'react';
import Background from './Background';

import '../CSS/page1.css';


function Page1(props)
{
    useEffect(()=> 
    {
        /* Ensure the model is shown */
        let model = document.getElementById("model");
        let navbar = document.getElementById("navbar");
        navbar.style.display = "block";
        model.style.display = "none";

        model.style.display = "none";
        /* animation for the search bar */
        let search = document.querySelector(".search-area");
        setTimeout(() =>
        {
            search.style.opacity = "1";
            search.style.animation = "appear 1.2s ease-in";
        }, 1400);

        /* animation for the pan elements */
        let pan = document.querySelectorAll(".pan");
        let pag = document.getElementById("pag");
        setTimeout(() =>
        {
            for(let i = 0; i < pan.length; i ++)
            {
                pan[i].style.display = "block";
                pan[i].style.animation = "appear 1.2s ease-in";
            }
            pag.style.display = "block";
            pag.style.animation = "appear 1.4s ease-in";
        }, 1400);

        /* filter element animation */
        let filters = document.querySelector(".filter").children;
        setTimeout(() =>
        {
            for(let i = 0; i < filters.length; i ++)
            {
                filters[i].style.display = "block";
                filters[i].style.animation = "appear 1.2s ease-in";
            }
        }, 1400);

        /* filter image script to show when clicked on */
        let filter = document.querySelectorAll(".filter-elements");
        let filter_img = document.querySelectorAll(".filter-img");
        let C_filter = document.querySelector(".filter-button");
        for(let i = 0; i < filter.length; i++)
        {
            filter[i].addEventListener("click", () =>
            {
                filter_img[i].style.display = "block";
                filter[i].style.backgroundColor = "rgba(64, 165, 24, 0.7)";
            });

            C_filter.addEventListener("click", () =>
            {
                filter_img[i].style.display = "none";
                filter[i].style.backgroundColor = "rgba(61, 61, 61, 0.7)";
            });
        }

        /* Hover brightens the color of the pan element details */
        let pan_details = document.querySelectorAll(".pan-details");
        let pan_price = document.querySelectorAll(".pan-price");

        for(let i = 0; i < pan.length; i++)
        {
            pan[i].onmouseover = function(event)
            {
                pan_details[i].style.color = "rgb(240, 248, 255, 0.8)";
                pan_price[i].style.color = "rgb(240, 248, 255, 0.8)";
            }
            pan[i].onmouseout = function(event)
            {
                pan_details[i].style.color = "black";
                pan_price[i].style.color = "black"; 
            }
        }

        /* animation for the pagination elements 
        let num = document.querySelector(".pagination").childNodes;
        for(let i = 1; i < num.length - 1; i++)
        {
            num[i].addEventListener("click", () =>
            {
                let prev = document.querySelector(".activee");
                prev.className = "";
                num[i].className = "activee";
            });
        }
        */



        /* Script to automatically format the number of elements on each page */
        const content = document.querySelector('.main-elements'); 
        const itemsPerPage = 4;
        let currentPage = 0;
        const items = Array.from(content.getElementsByClassName('pan'));
        console.log("YE");
        console.log(items);
          
        function showPage(page) 
        {
        const startIndex = page * itemsPerPage;
        const endIndex = startIndex + itemsPerPage;
        items.forEach((item, index) => 
        {
            item.classList.toggle('hidden', index < startIndex || index >= endIndex);
        });
        updateActiveButtonStates();
        }
        
        function createPageButtons() 
        {
        const totalPages = Math.ceil(items.length / itemsPerPage);
        const paginationContainer = document.createElement('div');
        const paginationDiv = document.body.appendChild(paginationContainer);
        paginationContainer.classList.add('pagination');
        
        
        // Add page buttons
            for (let i = 0; i < totalPages; i++) 
            {
                const pageButton = document.createElement('button');
                pageButton.textContent = i + 1;
                pageButton.addEventListener('click', () => 
                {
                currentPage = i;
                showPage(currentPage);
                updateActiveButtonStates();
                });
            
                content.appendChild(paginationContainer);
                paginationDiv.appendChild(pageButton);
                pag.appendChild(paginationDiv);
            }
        }
        
        function updateActiveButtonStates() 
        {
        const pageButtons = document.querySelectorAll('.pagination button');
        pageButtons.forEach((button, index) => 
        {
            if (index === currentPage) 
            {
            button.classList.add('active');
            } 
            else 
            {
            button.classList.remove('active');
            }
        });
        }
        
        createPageButtons(); // Call this function to create the page buttons initially
        showPage(currentPage);
    
    }, []);

    return (
        <>
            <Background />
            <div className = "filter" style = {{display: props.filter_display}}>
                <div className = "filter-title"><b>Available Filters:</b></div>
                <br /><br />
                <div className = "filter-elements">Filter 1<div className = "filter-img"></div></div>
                <div className = "filter-elements">Filter 2<div className = "filter-img"></div></div>
                <div className = "filter-elements">Filter 3<div className = "filter-img"></div></div>
                <div className = "filter-elements">Filter 4<div className = "filter-img"></div></div>
                <br />

                <button className = "filter-button">Clear Filter</button>
            </div>
            <div className = "main" style = {{display: props.main_display, backgroundColor: props.main_bgc, top: props.main_top,
            left: props.main_left, transform: props.transform, width: props.width, height: props.height, animation: props.animation}}>
                <div className = "search">
                    <form className = "search-area">
                        <input className ="search-area" type="search" placeholder="Search..." />
                    </form>    
                </div>
                <div className = "main-elements">
                    <div className = "pan">
                        <div className = "pan-img"></div>
                        <div className = "pan-details">
                            <a className = "p-d-title">Product Title</a>
                            <br/><br/>

                            <a className = "p-d-code">Product code</a>
                            <br/><br/>

                            <a className = "p-d-options">Options</a> | <a className = "p-d-category">Category</a> | <a className = "p-d-type">Type</a> | <a className = "p-d-vendor">Vendor</a>
                        </div>
                        <div className = "pan-price">
                            R1200 - R1400
                        </div>
                    </div>
                    <div className = "pan">
                        <div className = "pan-img"></div>
                        <div className = "pan-details">
                            <a className = "p-d-title">Product Title</a>
                            <br/><br/>

                            <a className = "p-d-code">Product code</a>
                            <br/><br/>

                            <a className = "p-d-options">Options</a> | <a className = "p-d-category">Category</a> | <a className = "p-d-type">Type</a> | <a className = "p-d-vendor">Vendor</a>
                        </div>
                        <div className = "pan-price">
                            R1200 - R1400
                        </div>
                    </div>
                    <div className = "pan">
                        <div className = "pan-img"></div>
                        <div className = "pan-details">
                            <a className = "p-d-title">Product Title</a>
                            <br/><br/>

                            <a className = "p-d-code">Product code</a>
                            <br/><br/>

                            <a className = "p-d-options">Options</a> | <a className = "p-d-category">Category</a> | <a className = "p-d-type">Type</a> | <a className = "p-d-vendor">Vendor</a>
                        </div>
                        <div className = "pan-price">
                            R1200 - R1400
                        </div>
                    </div>
                    <div className = "pan">
                        <div className = "pan-img"></div>
                        <div className = "pan-details">
                            <a className = "p-d-title">Product Title</a>
                            <br/><br/>

                            <a className = "p-d-code">Product code</a>
                            <br/><br/>

                            <a className = "p-d-options">Options</a> | <a className = "p-d-category">Category</a> | <a className = "p-d-type">Type</a> | <a className = "p-d-vendor">Vendor</a>
                        </div>
                        <div className = "pan-price">
                            R1200 - R1400
                        </div>
                    </div>

                    <div className = "pan">
                        <div className = "pan-img"></div>
                        <div className = "pan-details">
                            <a className = "p-d-title">Product Title</a>
                            <br/><br/>

                            <a className = "p-d-code">Product code</a>
                            <br/><br/>

                            <a className = "p-d-options">Options</a> | <a className = "p-d-category">Category</a> | <a className = "p-d-type">Type</a> | <a className = "p-d-vendor">Vendor</a>
                        </div>
                        <div className = "pan-price">
                            R1200 - R1400
                        </div>
                    </div>

                    <div className = "pan">
                        <div className = "pan-img"></div>
                        <div className = "pan-details">
                            <a className = "p-d-title">Product Title</a>
                            <br/><br/>

                            <a className = "p-d-code">Product code</a>
                            <br/><br/>

                            <a className = "p-d-options">Options</a> | <a className = "p-d-category">Category</a> | <a className = "p-d-type">Type</a> | <a className = "p-d-vendor">Vendor</a>
                        </div>
                        <div className = "pan-price">
                            R1200 - R1400
                        </div>
                    </div>

                    <div className = "pan">
                        <div className = "pan-img"></div>
                        <div className = "pan-details">
                            <a className = "p-d-title">Product Title</a>
                            <br/><br/>

                            <a className = "p-d-code">Product code</a>
                            <br/><br/>

                            <a className = "p-d-options">Options</a> | <a className = "p-d-category">Category</a> | <a className = "p-d-type">Type</a> | <a className = "p-d-vendor">Vendor</a>
                        </div>
                        <div className = "pan-price">
                            R1200 - R1400
                        </div>
                    </div>
                </div>
                <div className = "center" id = "pag"></div>
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
    animation: 'SlideUp2 1.2s ease-in'
};

export default Page1;