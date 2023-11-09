import {useEffect, useState} from 'react';
import $ from 'jquery';
import Page1 from '../components/Page1';
import Pan_details from '../components/semi-components/pan-detail';
import '../CSS/page1.css';

function Products(props)
{
    const [data, setData] = useState([]);

    useEffect(()=> 
    {
        /* Ensures the navbar is set correctly */
        let navigation = document.getElementById("navbar");
        window.onload = function(event)
        {
            navigation.style.left = "30%";
            navigation.style.position = "absolute";
            navigation.style.width = "70%";
            navigation.style.animation = "MoveLeft 1.2s ease";
        }

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
        
        
    }, []);

    return (
        <>
            <div className = "main">
                <div className = "search">
                    <form className = "search-area">
                        <input className ="search-area" type="search" placeholder="Search..." />
                    </form>    
                </div>
                <div className = "main-elements">
                    {data.map((_data, id)=>
                        {
                            return <Pan_details />

                        })
                    }
                    <Pan_details />
                </div>
                <div className = "center" id = "pag"></div>
            </div>

            <Page1 />

        </>
    );
}

export default Products;