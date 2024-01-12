import {useEffect} from 'react';
import {useState} from "react";
import $ from 'jquery';
import '../../CSS/login.css';
function Import_Product()
{
    const [file, setFile] = useState();

    const fileReader = new FileReader();

    const handleOnChange = (e) => 
    {
        setFile(e.target.files[0]);
    };
    const handleOnSubmit = (e) => 
    {
        e.preventDefault();
        if (file) 
        {
            fileReader.onload = function (event) 
            {
                const csvOutput = event.target.result;
            };
            fileReader.readAsText(file);
            console.log(file);
            const formData = new FormData();
            formData.append('file', file);
            const api_key = localStorage.getItem('api_key');
            
            $.ajaxSetup({ headers: { 'Authorization': 'ApiKey ' + api_key}, processData: false, contentType: false});

            $.post("http://localhost:8080/api/products/import", formData, [], 'multipart/form-data')
            .done(function( _data) 
            {
                console.log(_data);
            })
            .fail( function(xhr) 
            {
                alert(xhr.responseText);
            });
        }
    };

    useEffect(() =>
    {
        /* Fix any incorrect elements */
        let navigation = document.getElementById("navbar");
        let modal = document.getElementById("model");
        navigation.style.zIndex = "3";
        modal.style.display = "block";
        navigation.style.animation = "MoveRight 1.2s ease";
        navigation.style.position = "fixed";
        navigation.style.left = "0%";
        navigation.style.width = "100%";

    }, []);

    return (
        <>
            <div className = 'modal1' id = "model" style={{zIndex: '2', background:'linear-gradient(to bottom, #202020c7, #111119f0)'}}>

                <form className = 'modal-content' style ={{backgroundColor: 'none'}} method = 'post' autoComplete='off' id = 'form1'>

                    <div style = {{position: 'relative', top: '40%'}}>

                        <label style = {{fontSize: '18px',color: 'white', textDecoration:'underline'}}>
                            <b>Import Product</b> 
                            <i className = "info_icon" 
                            title='Ensure that your imported CSV file has the relative columns and information, to avoid any errors'></i>
                        </label>
                        <br /><br /><br />
                        <label style = {{color: 'white'}}><b>Imports customized products to the application</b></label>
                        <br /><br /><br />
                        
                        <input style = {{color: 'white'}} type={"file"} id = "file-upload-button" name = "file" accept={".csv"} onChange={handleOnChange}/>
                        <br /><br />
                        <button className = "button" onClick={(e) => { handleOnSubmit(e); }}>IMPORT CSV</button>
                    </div>
                    
                </form>
            </div>  
            <div className = "output"></div>  
        </>
    );
};

export default Import_Product;