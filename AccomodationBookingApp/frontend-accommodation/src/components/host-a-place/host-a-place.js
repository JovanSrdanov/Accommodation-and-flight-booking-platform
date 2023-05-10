import React, {useState} from 'react';
import ImageList from '@mui/material/ImageList';
import ImageListItem from '@mui/material/ImageListItem';
import CancelIcon from '@mui/icons-material/Cancel';
import {Button} from "@mui/material";

function HostAPlace() {
    const [uploadedImages, setUploadedImages] = useState([]);

    const handleImageRemove = (index) => {
        const updatedItems = [...uploadedImages];
        updatedItems.splice(index, 1);
        setUploadedImages(updatedItems);
    };

    const handleImageUpload = (event) => {
        const files = Array.from(event.target.files);

        const uploadedItems = files.filter(file => file.type.startsWith('image/'))
            .map((file) => ({
                img: URL.createObjectURL(file),
                file,
            }));

        setUploadedImages([...uploadedImages, ...uploadedItems]);
    };


    return (
        <div className="wrapper">
            <ImageList
                sx={{width: 500, height: 450, border: '1px solid #f57c00'}}
                cols={3}
                rowHeight={164}
            >
                {uploadedImages.map((item, index) => (
                    <ImageListItem key={item.img}>
                        <img src={item.img} alt="" loading="lazy"/>
                        <button
                            className="remove-image-button"
                            onClick={() => handleImageRemove(index)}
                        >
                            <CancelIcon/>
                        </button>
                    </ImageListItem>
                ))}
            </ImageList>
            <Button variant="contained" color="warning" component="label">
                Add pictures
                <input
                    type="file"
                    onChange={handleImageUpload}
                    multiple
                    accept="image/*"
                    style={{display: 'none'}}
                />
            </Button>


        </div>
    );
}

export default HostAPlace;
