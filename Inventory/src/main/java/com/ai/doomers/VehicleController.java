package com.ai.doomers;


import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
public class VehicleController {

    private final VehicleService vehicleService;

    public VehicleController(VehicleService vehicleService) {
        this.vehicleService = vehicleService;
        vehicleService.inventoryProcessor();
    }

    @GetMapping("/find/cars")
    public List<Vehicle> findCars(@RequestParam int capacity,  @RequestParam(required = false) String flags) {
        System.out.println("Capacity: " + capacity);
        System.out.println("Flags String: " + flags);
        return vehicleService.findVehicles(capacity, VehicleFlag.parseCommaSeperated(flags));
    }

    @PostMapping("/mark/in/use")
    public String markVehicleInUse(@RequestParam int id) {
        vehicleService.markVehicleAsUsed(id);
        System.out.println("Marked: " + id);

        return "";
    }

    @PostMapping("/mark/ride/over")
    public String returningVehicles(@RequestParam int id ){
        vehicleService.returningVehicles(id);

        return "";   
    }
}
