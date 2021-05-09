import { Card,Box,Divider, Button } from '@material-ui/core';
import React from 'react';
import { makeStyles } from '@material-ui/core/styles';
import Typography from '@material-ui/core/Typography';

const useStyles = makeStyles((theme) => ({
    root: {
        minWidth: 225,
        maxWidth: 225,
        padding: 5,
        margin: 1,
        borderRadius: 5,
        justifyContent: "space-around",
    },
    topBar: {
        height:3,
    },
    header: {
        fontSize: 30,
        fontFamily: 'Segoe UI Symbol',
        padding: 5,
    },
    cont: {
        fontSize: 20,
        fontFamily: 'Segoe UI Symbol',
    },
    imageBox: {
        borderRadius: 200,
    },
    logo:{
        height: 100,
        center:50,
    },
    titleDiv:{
        // backgroundColor:"#fff",
        padding: 1,
        boxShadow: 0,
    },
    priceDiv:{
        backgroundColor:"#ecf0f1",
        padding: 1,
        boxShadow: 0,
    },
    speedDiv:{
        // backgroundColor:"#faffb5",
        padding: 1,
    },
    durationDiv:{
        // backgroundColor:"#f6dcfa",
    },
    description:{
        fontSize:10,
        fontFamily: 'Segoe UI Symbol',
    },
}) )

const PlanCard = ( plan ) => {
    const classes = useStyles();

    return (
        <div>
            <Card className={classes.root} >
                {/* <Box className={classes.imageBox} >
                    <center>
                        <img className={classes.logo} src="plan.png" />
                    </center>
                </Box> */}
                <Box className={classes.topBar} bgcolor="secondary.main" color="white" ></Box>

                <Box className={classes.titleDiv} >
                    <center>
                    <Typography className={classes.header}>
                        { plan.planName }
                    </Typography>
                    <Typography className={classes.description} >
                        {plan.planDescription}
                    </Typography>
                    </center>
                </Box>
                <Divider />
                <Box className={classes.priceDiv} >
                    <center>
                    <Typography>
                        Price
                    </Typography>
                    <Typography className={classes.header}>
                        {plan.planPrice}
                    </Typography>
                    </center>
                </Box>

                <Box boxShadow={0} className={classes.speedDiv} >
                    <center>
                        <Typography>
                            Speed
                        </Typography>
                        <Typography className={classes.header}>
                        { plan.planSpeed }
                        </Typography>
                    </center>
                </Box>

                <Box boxShadow={0} className={classes.durationDiv} >
                    <center>
                        <Typography>
                            Duration
                        </Typography>
                        <Typography className={classes.header}>
                        { plan.planDuration }
                        </Typography>
                    </center>
                </Box>
                <center>
                <Button style={{margin:5}} variant="outlined" color="primary">Subscribe</Button>
                <p style={{fontSize:10, color:"grey"}} >Terms and Conditions applied</p></center>
                <Box className={classes.topBar} bgcolor="secondary.main" color="white" ></Box>
            </Card>
        </div>
    )
}
export default PlanCard
