package com.ai.doomers;


import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
public class VehicleController {

    private final VehicleService vehicleService;

    public VehicleController(VehicleService vehicleService) {
        this.vehicleService = vehicleService;
    }

    @GetMapping("/find/cars")
    public List<Vehicle> findCars(@RequestParam int capacity, 
                                  @RequestParam boolean wheelchairADA, 
                                  @RequestParam boolean visionImpairedADA) {
        return vehicleService.findVehicles(capacity, wheelchairADA, visionImpairedADA);
    }

    @PostMapping("/mark/in/use")
    public void markVehicleInUse(@RequestParam int id) {
        vehicleService.markVehicleAsUsed(id);
    }

    @PostMapping("/mark/ride/over")
    public void returningVehicles(@RequestParam int id ){
        vehicleService.returningVehicles(id);       
    }
}
