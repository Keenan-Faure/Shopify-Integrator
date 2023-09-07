import {useEffect} from 'react';
import Background from './Background';

import '../CSS/page1.css';


function Page1(props)
{
    useEffect(()=> 
    {

        /* Ensures the navbar is set correctly */
        let navigation = document.getElementById("navbar");
        let search = document.querySelector(".search-area");
        window.onload = function(event)
        {
            navigation.style.left = "30%";
            navigation.style.position = "absolute";
            navigation.style.width = "70%";
            navigation.style.animation = "MoveLeft 1.2s ease";
            setTimeout(() =>
            {
                search.style.opacity = "1";
                search.style.animation = "appear 1.2s ease-in";
            }, 1400);
        }

        /* animation for the pan elements */
        let pan = document.querySelectorAll(".pan");
        setTimeout(() =>
        {
            for(let i = 0; i < pan.length; i ++)
            {
                pan[i].style.display = "block";
                pan[i].style.animation = "appear 1.2s ease-in";
            }
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
            <div className = "filter">
                <div className = "filter-title"><b>Available Filters:</b></div>
                <br /><br />
                <div className = "filter-elements">Filter 1<div className = "filter-img"></div></div>
                <div className = "filter-elements">Filter 2<div className = "filter-img"></div></div>
                <div className = "filter-elements">Filter 3<div className = "filter-img"></div></div>
                <div className = "filter-elements">Filter 4<div className = "filter-img"></div></div>
                <br />

                <button className = "filter-button">Clear Filter</button>
            </div>
            <div className = "main">
                <div className = "search">
                    <form className = "search-area">
                        <input className ="search-area" type="search" placeholder="Search..." />
                    </form>    
                </div>
                <div className = "main-elements">
                    <div className = "pan"></div>
                    <div className = "pan"></div>
                    <div className = "pan"></div>
                    <div className = "pan"></div>
                </div>
            </div>
        </>
    );
}

export default Page1;