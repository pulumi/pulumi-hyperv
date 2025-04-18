// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.hyperv.machine.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import com.pulumi.exceptions.MissingRequiredPropertyException;
import java.lang.Integer;
import java.lang.String;
import java.util.Objects;


public final class HardDriveInputArgs extends com.pulumi.resources.ResourceArgs {

    public static final HardDriveInputArgs Empty = new HardDriveInputArgs();

    @Import(name="controllerLocation", required=true)
    private Output<Integer> controllerLocation;

    public Output<Integer> controllerLocation() {
        return this.controllerLocation;
    }

    @Import(name="controllerNumber", required=true)
    private Output<Integer> controllerNumber;

    public Output<Integer> controllerNumber() {
        return this.controllerNumber;
    }

    @Import(name="controllerType", required=true)
    private Output<String> controllerType;

    public Output<String> controllerType() {
        return this.controllerType;
    }

    @Import(name="path", required=true)
    private Output<String> path;

    public Output<String> path() {
        return this.path;
    }

    private HardDriveInputArgs() {}

    private HardDriveInputArgs(HardDriveInputArgs $) {
        this.controllerLocation = $.controllerLocation;
        this.controllerNumber = $.controllerNumber;
        this.controllerType = $.controllerType;
        this.path = $.path;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(HardDriveInputArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private HardDriveInputArgs $;

        public Builder() {
            $ = new HardDriveInputArgs();
        }

        public Builder(HardDriveInputArgs defaults) {
            $ = new HardDriveInputArgs(Objects.requireNonNull(defaults));
        }

        public Builder controllerLocation(Output<Integer> controllerLocation) {
            $.controllerLocation = controllerLocation;
            return this;
        }

        public Builder controllerLocation(Integer controllerLocation) {
            return controllerLocation(Output.of(controllerLocation));
        }

        public Builder controllerNumber(Output<Integer> controllerNumber) {
            $.controllerNumber = controllerNumber;
            return this;
        }

        public Builder controllerNumber(Integer controllerNumber) {
            return controllerNumber(Output.of(controllerNumber));
        }

        public Builder controllerType(Output<String> controllerType) {
            $.controllerType = controllerType;
            return this;
        }

        public Builder controllerType(String controllerType) {
            return controllerType(Output.of(controllerType));
        }

        public Builder path(Output<String> path) {
            $.path = path;
            return this;
        }

        public Builder path(String path) {
            return path(Output.of(path));
        }

        public HardDriveInputArgs build() {
            if ($.controllerLocation == null) {
                throw new MissingRequiredPropertyException("HardDriveInputArgs", "controllerLocation");
            }
            if ($.controllerNumber == null) {
                throw new MissingRequiredPropertyException("HardDriveInputArgs", "controllerNumber");
            }
            if ($.controllerType == null) {
                throw new MissingRequiredPropertyException("HardDriveInputArgs", "controllerType");
            }
            if ($.path == null) {
                throw new MissingRequiredPropertyException("HardDriveInputArgs", "path");
            }
            return $;
        }
    }

}
