import {useEffect} from 'react';
import {useState} from "react";
import $ from 'jquery';
import '../../CSS/login.css';
import Background from '../../components/Background';

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
            
            $.ajaxSetup({ headers: { 'Authorization': 'ApiKey ' + api_key}, processData: false, contentType: false, method: 'post',});

            $.post("http://localhost:8080/api/products/import?test=true", formData, [], 'multipart/form-data')
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
        modal.style.display = "block";
        navigation.style.animation = "MoveRight 1.2s ease";
        navigation.style.position = "fixed";
        navigation.style.left = "0%";
        navigation.style.width = "100%";

    }, []);

    return (
        <>
            <Background />
            <div className = 'modal1' id = "model" style={{zIndex: '2', background:'linear-gradient(to bottom, #202020c7, #111119f0)'}}>

                <form className = 'modal-content' style ={{backgroundColor: 'none'}} method = 'post' autoComplete='off' id = 'form1'>

                    <div style = {{position: 'relative', top: '40%'}}>

                        <input type={"file"} id = "file-upload-button" name = "_import" accept={".csv"} onChange={handleOnChange}/>
                        <br /><br />
                        <button className = "button" onClick={(e) => { handleOnSubmit(e); }}> IMPORT CSV </button>
                    </div>
                    
                </form>
            </div>  
            <div className = "output"></div>  
        </>
    );
};
  
Import_Product.defaultProps = 
{

};
export default Import_Product;