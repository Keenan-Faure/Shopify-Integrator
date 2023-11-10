import {useEffect, useState} from 'react';
import $ from 'jquery';
import Page1 from '../components/Page1';
import Pan_details from '../components/semi-components/pan-detail';
import '../CSS/page1.css';
import product from '../media/products.png';

function Products(props)
{
    const [data, setData] = useState([]);

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
        
        
    }, []);

    return (
        <div className = "products">
            <div className = "main">
                <div className = "search">
                    <form className = "search-area">
                        <input className ="search-area" type="search" placeholder="Search..." />
                    </form>    
                </div>
                <div className = "main-elements">
                    <div className = "pan-main">
                        {data.map((_data, id)=>
                            {
                                return <Pan_details />

                            })
                        }
                        <Pan_details />
                    </div>
                </div>
                <div className = "center" id = "pag">
                    
                </div>
            </div>

            <Page1 image = {product} title = "Products"/>

        </div>
    );
}

export default Products;