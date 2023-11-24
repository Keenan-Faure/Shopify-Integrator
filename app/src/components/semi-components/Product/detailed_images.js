import {useEffect} from 'react';
import '../../../CSS/detailed.css';

function Detailed_Images(props)
{
    useEffect(()=> 
    {
        let slideIndex = 0;
        showSlides();
        
        function showSlides() 
        {
            let i;
            let slides = document.getElementsByClassName("mySlides");
            console.log(slides);
            if(slides != null)
            {
                for (i = 0; i < slides.length; i++) 
                {
                slides[i].style.display = "none";  
                }
                slideIndex++;
                if (slideIndex > slides.length) {slideIndex = 1}    
                slides[slideIndex-1].style.display = "block";  
                setTimeout(() =>
                {
                    showSlides();
                }, 5000); // Change image every 5 seconds
            }
            else 
            {
                console.log("no slides to display");
            }
            
        }
    }, []);

    return (
    <>
        <div className = "mySlides fade">
            <img src = {props.Image1} className = "details-image"></img>
        </div>
    </>
    );
  
};

Detailed_Images.defaultProps = 
{ 
    Image1: '#ccc',
}
export default Detailed_Images;