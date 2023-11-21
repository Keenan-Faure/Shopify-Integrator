import {useEffect} from 'react';
import Background from './Background';
import $ from 'jquery';
import '../CSS/page1.css';


function Page1(props)
{
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
    }, []);

    return (
        <>
            <Background />
            <div className = "filter" style = {{display: props.filter_display}}>
                <div className = "filter-title"><b>Available Filters:</b></div>
                <div className = "filter-elements">Filter 1<div className = "filter-img"></div></div>
                <div className = "filter-elements">Filter 2<div className = "filter-img"></div></div>
                <div className = "filter-elements">Filter 3<div className = "filter-img"></div></div>
                <div className = "filter-elements">Filter 4<div className = "filter-img"></div></div>
                <br />

                <button className = "filter-button">Clear Filter</button>
                
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