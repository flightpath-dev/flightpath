package common

import "fmt"

// MavMessageId
// Represents a MAVLink message ID from the common dialect.
// Reference: https://mavlink.io/en/messages/common.html
type MavMessageId uint32

const (
	MavMessageIdHeartbeat                          MavMessageId = 0
	MavMessageIdSysStatus                          MavMessageId = 1
	MavMessageIdSystemTime                         MavMessageId = 2
	MavMessageIdPing                               MavMessageId = 4
	MavMessageIdChangeOperatorControl              MavMessageId = 5
	MavMessageIdChangeOperatorControlAck           MavMessageId = 6
	MavMessageIdAuthKey                            MavMessageId = 7
	MavMessageIdLinkNodeStatus                     MavMessageId = 8
	MavMessageIdSetMode                            MavMessageId = 11
	MavMessageIdParamRequestRead                   MavMessageId = 20
	MavMessageIdParamRequestList                   MavMessageId = 21
	MavMessageIdParamValue                         MavMessageId = 22
	MavMessageIdParamSet                           MavMessageId = 23
	MavMessageIdGpsRawInt                          MavMessageId = 24
	MavMessageIdGpsStatus                          MavMessageId = 25
	MavMessageIdScaledImu                          MavMessageId = 26
	MavMessageIdRawImu                             MavMessageId = 27
	MavMessageIdRawPressure                        MavMessageId = 28
	MavMessageIdScaledPressure                     MavMessageId = 29
	MavMessageIdAttitude                           MavMessageId = 30
	MavMessageIdAttitudeQuaternion                 MavMessageId = 31
	MavMessageIdLocalPositionNed                   MavMessageId = 32
	MavMessageIdGlobalPositionInt                  MavMessageId = 33
	MavMessageIdRcChannelsScaled                   MavMessageId = 34
	MavMessageIdRcChannelsRaw                      MavMessageId = 35
	MavMessageIdServoOutputRaw                     MavMessageId = 36
	MavMessageIdMissionRequestPartialList          MavMessageId = 37
	MavMessageIdMissionWritePartialList            MavMessageId = 38
	MavMessageIdMissionItem                        MavMessageId = 39
	MavMessageIdMissionRequest                     MavMessageId = 40
	MavMessageIdMissionSetCurrent                  MavMessageId = 41
	MavMessageIdMissionCurrent                     MavMessageId = 42
	MavMessageIdMissionRequestList                 MavMessageId = 43
	MavMessageIdMissionCount                       MavMessageId = 44
	MavMessageIdMissionClearAll                    MavMessageId = 45
	MavMessageIdMissionItemReached                 MavMessageId = 46
	MavMessageIdMissionAck                         MavMessageId = 47
	MavMessageIdSetGpsGlobalOrigin                 MavMessageId = 48
	MavMessageIdGpsGlobalOrigin                    MavMessageId = 49
	MavMessageIdParamMapRc                         MavMessageId = 50
	MavMessageIdMissionRequestInt                  MavMessageId = 51
	MavMessageIdSafetySetAllowedArea               MavMessageId = 54
	MavMessageIdSafetyAllowedArea                  MavMessageId = 55
	MavMessageIdAttitudeQuaternionCov              MavMessageId = 61
	MavMessageIdNavControllerOutput                MavMessageId = 62
	MavMessageIdGlobalPositionIntCov               MavMessageId = 63
	MavMessageIdLocalPositionNedCov                MavMessageId = 64
	MavMessageIdRcChannels                         MavMessageId = 65
	MavMessageIdRequestDataStream                  MavMessageId = 66
	MavMessageIdDataStream                         MavMessageId = 67
	MavMessageIdManualControl                      MavMessageId = 69
	MavMessageIdRcChannelsOverride                 MavMessageId = 70
	MavMessageIdMissionItemInt                     MavMessageId = 73
	MavMessageIdVfrHud                             MavMessageId = 74
	MavMessageIdCommandInt                         MavMessageId = 75
	MavMessageIdCommandLong                        MavMessageId = 76
	MavMessageIdCommandAck                         MavMessageId = 77
	MavMessageIdCommandCancel                      MavMessageId = 80
	MavMessageIdManualSetpoint                     MavMessageId = 81
	MavMessageIdSetAttitudeTarget                  MavMessageId = 82
	MavMessageIdAttitudeTarget                     MavMessageId = 83
	MavMessageIdSetPositionTargetLocalNed          MavMessageId = 84
	MavMessageIdPositionTargetLocalNed             MavMessageId = 85
	MavMessageIdSetPositionTargetGlobalInt         MavMessageId = 86
	MavMessageIdPositionTargetGlobalInt            MavMessageId = 87
	MavMessageIdPositionTargetGlobalIntRelHome     MavMessageId = 88
	MavMessageIdLocalPositionNedSystemGlobalOffset MavMessageId = 89
	MavMessageIdHilState                           MavMessageId = 90
	MavMessageIdHilControls                        MavMessageId = 91
	MavMessageIdHilRcInputsRaw                     MavMessageId = 92
	MavMessageIdHilActuatorControls                MavMessageId = 93
	MavMessageIdOpticalFlow                        MavMessageId = 100
	MavMessageIdGlobalVisionPositionEstimate       MavMessageId = 101
	MavMessageIdVisionPositionEstimate             MavMessageId = 102
	MavMessageIdVisionSpeedEstimate                MavMessageId = 103
	MavMessageIdViconPositionEstimate              MavMessageId = 104
	MavMessageIdHighresImu                         MavMessageId = 105
	MavMessageIdOpticalFlowRad                     MavMessageId = 106
	MavMessageIdHilSensor                          MavMessageId = 107
	MavMessageIdSimState                           MavMessageId = 108
	MavMessageIdRadioStatus                        MavMessageId = 109
	MavMessageIdFileTransferProtocol               MavMessageId = 110
	MavMessageIdTimesync                           MavMessageId = 111
	MavMessageIdCameraTrigger                      MavMessageId = 112
	MavMessageIdHilGps                             MavMessageId = 113
	MavMessageIdHilOpticalFlow                     MavMessageId = 114
	MavMessageIdHilStateQuaternion                 MavMessageId = 115
	MavMessageIdScaledImu2                         MavMessageId = 116
	MavMessageIdLogRequestList                     MavMessageId = 117
	MavMessageIdLogEntry                           MavMessageId = 118
	MavMessageIdLogRequestData                     MavMessageId = 119
	MavMessageIdLogData                            MavMessageId = 120
	MavMessageIdLogErase                           MavMessageId = 121
	MavMessageIdLogRequestEnd                      MavMessageId = 122
	MavMessageIdGpsInjectData                      MavMessageId = 123
	MavMessageIdGps2Raw                            MavMessageId = 124
	MavMessageIdPowerStatus                        MavMessageId = 125
	MavMessageIdSerialControl                      MavMessageId = 126
	MavMessageIdGpsRtk                             MavMessageId = 127
	MavMessageIdGps2Rtk                            MavMessageId = 128
	MavMessageIdScaledImu3                         MavMessageId = 129
	MavMessageIdDataTransmissionHandshake          MavMessageId = 130
	MavMessageIdEncapsulatedData                   MavMessageId = 131
	MavMessageIdDistanceSensor                     MavMessageId = 132
	MavMessageIdTerrainRequest                     MavMessageId = 133
	MavMessageIdTerrainData                        MavMessageId = 134
	MavMessageIdTerrainCheck                       MavMessageId = 135
	MavMessageIdTerrainReport                      MavMessageId = 136
	MavMessageIdScaledPressure2                    MavMessageId = 137
	MavMessageIdAttPosMocap                        MavMessageId = 138
	MavMessageIdSetActuatorControlTarget           MavMessageId = 139
	MavMessageIdActuatorControlTarget              MavMessageId = 140
	MavMessageIdAltitude                           MavMessageId = 141
	MavMessageIdResourceRequest                    MavMessageId = 142
	MavMessageIdScaledPressure3                    MavMessageId = 143
	MavMessageIdFollowTarget                       MavMessageId = 144
	MavMessageIdControlSystemState                 MavMessageId = 146
	MavMessageIdBatteryStatus                      MavMessageId = 147
	MavMessageIdAutopilotVersion                   MavMessageId = 148
	MavMessageIdLandingTarget                      MavMessageId = 149
	MavMessageIdSensorOffsets                      MavMessageId = 150
	MavMessageIdSetMagOffsets                      MavMessageId = 151
	MavMessageIdMeminfo                            MavMessageId = 152
	MavMessageIdApAdc                              MavMessageId = 153
	MavMessageIdDigicamConfigure                   MavMessageId = 154
	MavMessageIdDigicamControl                     MavMessageId = 155
	MavMessageIdMountConfigure                     MavMessageId = 156
	MavMessageIdMountControl                       MavMessageId = 157
	MavMessageIdMountStatus                        MavMessageId = 158
	MavMessageIdFencePoint                         MavMessageId = 160
	MavMessageIdFenceFetchPoint                    MavMessageId = 161
	MavMessageIdFenceStatus                        MavMessageId = 162
	MavMessageIdAhrs                               MavMessageId = 163
	MavMessageIdSimstate                           MavMessageId = 164
	MavMessageIdHwstatus                           MavMessageId = 165
	MavMessageIdRadio                              MavMessageId = 166
	MavMessageIdLimitsStatus                       MavMessageId = 167
	MavMessageIdWind                               MavMessageId = 168
	MavMessageIdData16                             MavMessageId = 169
	MavMessageIdData32                             MavMessageId = 170
	MavMessageIdData64                             MavMessageId = 171
	MavMessageIdData96                             MavMessageId = 172
	MavMessageIdRangefinder                        MavMessageId = 173
	MavMessageIdAirspeedAutocal                    MavMessageId = 174
	MavMessageIdRallyPoint                         MavMessageId = 175
	MavMessageIdRallyFetchPoint                    MavMessageId = 176
	MavMessageIdCompassCalibrationProgress         MavMessageId = 177
	MavMessageIdEkfStatusReport                    MavMessageId = 179
	MavMessageIdPidTuning                          MavMessageId = 180
	MavMessageIdDeepstall                          MavMessageId = 181
	MavMessageIdGimbalReport                       MavMessageId = 182
	MavMessageIdGimbalControl                      MavMessageId = 183
	MavMessageIdGimbalTorqueCmdReport              MavMessageId = 184
	MavMessageIdMagCalReport                       MavMessageId = 192
	MavMessageIdEfiStatus                          MavMessageId = 225
	MavMessageIdEstimatorStatus                    MavMessageId = 230
	MavMessageIdWindCov                            MavMessageId = 231
	MavMessageIdGpsInput                           MavMessageId = 232
	MavMessageIdGpsRtcmData                        MavMessageId = 233
	MavMessageIdHighLatency                        MavMessageId = 234
	MavMessageIdHighLatency2                       MavMessageId = 235
	MavMessageIdVibration                          MavMessageId = 241
	MavMessageIdHomePosition                       MavMessageId = 242
	MavMessageIdSetHomePosition                    MavMessageId = 243
	MavMessageIdMessageInterval                    MavMessageId = 244
	MavMessageIdExtendedSysState                   MavMessageId = 245
	MavMessageIdAdsbVehicle                        MavMessageId = 246
	MavMessageIdCollision                          MavMessageId = 247
	MavMessageIdV2Extension                        MavMessageId = 248
	MavMessageIdMemoryVect                         MavMessageId = 249
	MavMessageIdDebugVect                          MavMessageId = 250
	MavMessageIdNamedValueFloat                    MavMessageId = 251
	MavMessageIdNamedValueInt                      MavMessageId = 252
	MavMessageIdStatustext                         MavMessageId = 253
	MavMessageIdDebug                              MavMessageId = 254
	MavMessageIdSetupSigning                       MavMessageId = 256
	MavMessageIdButtonChange                       MavMessageId = 257
	MavMessageIdPlayTune                           MavMessageId = 258
	MavMessageIdCameraInformation                  MavMessageId = 259
	MavMessageIdCameraSettings                     MavMessageId = 260
	MavMessageIdStorageInformation                 MavMessageId = 261
	MavMessageIdCameraCaptureStatus                MavMessageId = 262
	MavMessageIdCameraImageCaptured                MavMessageId = 263
	MavMessageIdFlightInformation                  MavMessageId = 264
	MavMessageIdMountOrientation                   MavMessageId = 265
	MavMessageIdLoggingData                        MavMessageId = 266
	MavMessageIdLoggingDataAcked                   MavMessageId = 267
	MavMessageIdLoggingAck                         MavMessageId = 268
	MavMessageIdVideoStreamInformation             MavMessageId = 269
	MavMessageIdVideoStreamStatus                  MavMessageId = 270
	MavMessageIdCameraFovStatus                    MavMessageId = 271
	MavMessageIdCameraTrackingImageStatus          MavMessageId = 275
	MavMessageIdCameraTrackingGeoStatus            MavMessageId = 276
	MavMessageIdCameraThermalRange                 MavMessageId = 277
	MavMessageIdGimbalManagerInformation           MavMessageId = 280
	MavMessageIdGimbalManagerStatus                MavMessageId = 281
	MavMessageIdGimbalManagerSetAttitude           MavMessageId = 282
	MavMessageIdGimbalDeviceInformation            MavMessageId = 283
	MavMessageIdGimbalDeviceSetAttitude            MavMessageId = 284
	MavMessageIdGimbalDeviceAttitudeStatus         MavMessageId = 285
	MavMessageIdAutopilotStateForGimbalDevice      MavMessageId = 286
	MavMessageIdGimbalManagerSetPitchyaw           MavMessageId = 287
	MavMessageIdGimbalManagerSetManualControl      MavMessageId = 288
	MavMessageIdEscInfo                            MavMessageId = 290
	MavMessageIdEscStatus                          MavMessageId = 291
	MavMessageIdAirspeed                           MavMessageId = 295
	MavMessageIdWifiConfigAp                       MavMessageId = 299
	MavMessageIdProtocolVersion                    MavMessageId = 300
	MavMessageIdAisVessel                          MavMessageId = 301
	MavMessageIdUavcanNodeStatus                   MavMessageId = 310
	MavMessageIdUavcanNodeInfo                     MavMessageId = 311
	MavMessageIdParamExtRequestRead                MavMessageId = 320
	MavMessageIdParamExtRequestList                MavMessageId = 321
	MavMessageIdParamExtValue                      MavMessageId = 322
	MavMessageIdParamExtSet                        MavMessageId = 323
	MavMessageIdParamExtAck                        MavMessageId = 324
	MavMessageIdObstacleDistance                   MavMessageId = 330
	MavMessageIdOdometry                           MavMessageId = 331
	MavMessageIdTrajectoryRepresentationWaypoints  MavMessageId = 332
	MavMessageIdTrajectoryRepresentationBezier     MavMessageId = 333
	MavMessageIdCellularStatus                     MavMessageId = 334
	MavMessageIdIsbdLinkStatus                     MavMessageId = 335
	MavMessageIdCellularConfig                     MavMessageId = 336
	MavMessageIdRawRpm                             MavMessageId = 339
	MavMessageIdUtmGlobalPosition                  MavMessageId = 340
	MavMessageIdParamError                         MavMessageId = 345
	MavMessageIdDebugFloatArray                    MavMessageId = 350
	MavMessageIdOrbitExecutionStatus               MavMessageId = 360
	MavMessageIdFigureEightExecutionStatus         MavMessageId = 361
	MavMessageIdSmartBatteryInfo                   MavMessageId = 370
	MavMessageIdFuelStatus                         MavMessageId = 371
	MavMessageIdBatteryInfo                        MavMessageId = 372
	MavMessageIdGeneratorStatus                    MavMessageId = 373
	MavMessageIdActuatorOutputStatus               MavMessageId = 375
	MavMessageIdTimeEstimateToTarget               MavMessageId = 380
	MavMessageIdTunnel                             MavMessageId = 385
	MavMessageIdCanFrame                           MavMessageId = 386
	MavMessageIdCanfdFrame                         MavMessageId = 387
	MavMessageIdCanFilterModify                    MavMessageId = 388
	MavMessageIdOnboardComputerStatus              MavMessageId = 390
	MavMessageIdComponentInformation               MavMessageId = 395
	MavMessageIdComponentInformationBasic          MavMessageId = 396
	MavMessageIdComponentMetadata                  MavMessageId = 397
	MavMessageIdComponentMetadataV2                MavMessageId = 398
	MavMessageIdPlayTuneV2                         MavMessageId = 400
	MavMessageIdSupportedTunes                     MavMessageId = 401
	MavMessageIdEvent                              MavMessageId = 410
	MavMessageIdCurrentEventSequence               MavMessageId = 411
	MavMessageIdRequestEvent                       MavMessageId = 412
	MavMessageIdResponseEventError                 MavMessageId = 413
	MavMessageIdAvailableModes                     MavMessageId = 435
	MavMessageIdCurrentMode                        MavMessageId = 436
	MavMessageIdAvailableModesMonitor              MavMessageId = 437
	MavMessageIdIlluminatorStatus                  MavMessageId = 440
	MavMessageIdWheelDistance                      MavMessageId = 9000
	MavMessageIdWinchStatus                        MavMessageId = 9005
	MavMessageIdOpenDroneIdBasicId                 MavMessageId = 12900
	MavMessageIdOpenDroneIdLocation                MavMessageId = 12901
	MavMessageIdOpenDroneIdAuthentication          MavMessageId = 12902
	MavMessageIdOpenDroneIdSelfId                  MavMessageId = 12903
	MavMessageIdOpenDroneIdSystem                  MavMessageId = 12904
	MavMessageIdOpenDroneIdOperatorId              MavMessageId = 12905
	MavMessageIdOpenDroneIdMessagePack             MavMessageId = 12915
	MavMessageIdOpenDroneIdArmStatus               MavMessageId = 12918
	MavMessageIdOpenDroneIdSystemUpdate            MavMessageId = 12919
	MavMessageIdHygrometerSensor                   MavMessageId = 12920
)

var mavMessageIdStrings = map[MavMessageId]string{
	MavMessageIdHeartbeat:                          "HEARTBEAT",
	MavMessageIdSysStatus:                          "SYS_STATUS",
	MavMessageIdSystemTime:                         "SYSTEM_TIME",
	MavMessageIdPing:                               "PING",
	MavMessageIdChangeOperatorControl:              "CHANGE_OPERATOR_CONTROL",
	MavMessageIdChangeOperatorControlAck:           "CHANGE_OPERATOR_CONTROL_ACK",
	MavMessageIdAuthKey:                            "AUTH_KEY",
	MavMessageIdLinkNodeStatus:                     "LINK_NODE_STATUS",
	MavMessageIdSetMode:                            "SET_MODE",
	MavMessageIdParamRequestRead:                   "PARAM_REQUEST_READ",
	MavMessageIdParamRequestList:                   "PARAM_REQUEST_LIST",
	MavMessageIdParamValue:                         "PARAM_VALUE",
	MavMessageIdParamSet:                           "PARAM_SET",
	MavMessageIdGpsRawInt:                          "GPS_RAW_INT",
	MavMessageIdGpsStatus:                          "GPS_STATUS",
	MavMessageIdScaledImu:                          "SCALED_IMU",
	MavMessageIdRawImu:                             "RAW_IMU",
	MavMessageIdRawPressure:                        "RAW_PRESSURE",
	MavMessageIdScaledPressure:                     "SCALED_PRESSURE",
	MavMessageIdAttitude:                           "ATTITUDE",
	MavMessageIdAttitudeQuaternion:                 "ATTITUDE_QUATERNION",
	MavMessageIdLocalPositionNed:                   "LOCAL_POSITION_NED",
	MavMessageIdGlobalPositionInt:                  "GLOBAL_POSITION_INT",
	MavMessageIdRcChannelsScaled:                   "RC_CHANNELS_SCALED",
	MavMessageIdRcChannelsRaw:                      "RC_CHANNELS_RAW",
	MavMessageIdServoOutputRaw:                     "SERVO_OUTPUT_RAW",
	MavMessageIdMissionRequestPartialList:          "MISSION_REQUEST_PARTIAL_LIST",
	MavMessageIdMissionWritePartialList:            "MISSION_WRITE_PARTIAL_LIST",
	MavMessageIdMissionItem:                        "MISSION_ITEM",
	MavMessageIdMissionRequest:                     "MISSION_REQUEST",
	MavMessageIdMissionSetCurrent:                  "MISSION_SET_CURRENT",
	MavMessageIdMissionCurrent:                     "MISSION_CURRENT",
	MavMessageIdMissionRequestList:                 "MISSION_REQUEST_LIST",
	MavMessageIdMissionCount:                       "MISSION_COUNT",
	MavMessageIdMissionClearAll:                    "MISSION_CLEAR_ALL",
	MavMessageIdMissionItemReached:                 "MISSION_ITEM_REACHED",
	MavMessageIdMissionAck:                         "MISSION_ACK",
	MavMessageIdSetGpsGlobalOrigin:                 "SET_GPS_GLOBAL_ORIGIN",
	MavMessageIdGpsGlobalOrigin:                    "GPS_GLOBAL_ORIGIN",
	MavMessageIdParamMapRc:                         "PARAM_MAP_RC",
	MavMessageIdMissionRequestInt:                  "MISSION_REQUEST_INT",
	MavMessageIdSafetySetAllowedArea:               "SAFETY_SET_ALLOWED_AREA",
	MavMessageIdSafetyAllowedArea:                  "SAFETY_ALLOWED_AREA",
	MavMessageIdAttitudeQuaternionCov:              "ATTITUDE_QUATERNION_COV",
	MavMessageIdNavControllerOutput:                "NAV_CONTROLLER_OUTPUT",
	MavMessageIdGlobalPositionIntCov:               "GLOBAL_POSITION_INT_COV",
	MavMessageIdLocalPositionNedCov:                "LOCAL_POSITION_NED_COV",
	MavMessageIdRcChannels:                         "RC_CHANNELS",
	MavMessageIdRequestDataStream:                  "REQUEST_DATA_STREAM",
	MavMessageIdDataStream:                         "DATA_STREAM",
	MavMessageIdManualControl:                      "MANUAL_CONTROL",
	MavMessageIdRcChannelsOverride:                 "RC_CHANNELS_OVERRIDE",
	MavMessageIdMissionItemInt:                     "MISSION_ITEM_INT",
	MavMessageIdVfrHud:                             "VFR_HUD",
	MavMessageIdCommandInt:                         "COMMAND_INT",
	MavMessageIdCommandLong:                        "COMMAND_LONG",
	MavMessageIdCommandAck:                         "COMMAND_ACK",
	MavMessageIdCommandCancel:                      "COMMAND_CANCEL",
	MavMessageIdManualSetpoint:                     "MANUAL_SETPOINT",
	MavMessageIdSetAttitudeTarget:                  "SET_ATTITUDE_TARGET",
	MavMessageIdAttitudeTarget:                     "ATTITUDE_TARGET",
	MavMessageIdSetPositionTargetLocalNed:          "SET_POSITION_TARGET_LOCAL_NED",
	MavMessageIdPositionTargetLocalNed:             "POSITION_TARGET_LOCAL_NED",
	MavMessageIdSetPositionTargetGlobalInt:         "SET_POSITION_TARGET_GLOBAL_INT",
	MavMessageIdPositionTargetGlobalInt:            "POSITION_TARGET_GLOBAL_INT",
	MavMessageIdPositionTargetGlobalIntRelHome:     "POSITION_TARGET_GLOBAL_INT_REL_HOME",
	MavMessageIdLocalPositionNedSystemGlobalOffset: "LOCAL_POSITION_NED_SYSTEM_GLOBAL_OFFSET",
	MavMessageIdHilState:                           "HIL_STATE",
	MavMessageIdHilControls:                        "HIL_CONTROLS",
	MavMessageIdHilRcInputsRaw:                     "HIL_RC_INPUTS_RAW",
	MavMessageIdHilActuatorControls:                "HIL_ACTUATOR_CONTROLS",
	MavMessageIdOpticalFlow:                        "OPTICAL_FLOW",
	MavMessageIdGlobalVisionPositionEstimate:       "GLOBAL_VISION_POSITION_ESTIMATE",
	MavMessageIdVisionPositionEstimate:             "VISION_POSITION_ESTIMATE",
	MavMessageIdVisionSpeedEstimate:                "VISION_SPEED_ESTIMATE",
	MavMessageIdViconPositionEstimate:              "VICON_POSITION_ESTIMATE",
	MavMessageIdHighresImu:                         "HIGHRES_IMU",
	MavMessageIdOpticalFlowRad:                     "OPTICAL_FLOW_RAD",
	MavMessageIdHilSensor:                          "HIL_SENSOR",
	MavMessageIdSimState:                           "SIM_STATE",
	MavMessageIdRadioStatus:                        "RADIO_STATUS",
	MavMessageIdFileTransferProtocol:               "FILE_TRANSFER_PROTOCOL",
	MavMessageIdTimesync:                           "TIMESYNC",
	MavMessageIdCameraTrigger:                      "CAMERA_TRIGGER",
	MavMessageIdHilGps:                             "HIL_GPS",
	MavMessageIdHilOpticalFlow:                     "HIL_OPTICAL_FLOW",
	MavMessageIdHilStateQuaternion:                 "HIL_STATE_QUATERNION",
	MavMessageIdScaledImu2:                         "SCALED_IMU2",
	MavMessageIdLogRequestList:                     "LOG_REQUEST_LIST",
	MavMessageIdLogEntry:                           "LOG_ENTRY",
	MavMessageIdLogRequestData:                     "LOG_REQUEST_DATA",
	MavMessageIdLogData:                            "LOG_DATA",
	MavMessageIdLogErase:                           "LOG_ERASE",
	MavMessageIdLogRequestEnd:                      "LOG_REQUEST_END",
	MavMessageIdGpsInjectData:                      "GPS_INJECT_DATA",
	MavMessageIdGps2Raw:                            "GPS2_RAW",
	MavMessageIdPowerStatus:                        "POWER_STATUS",
	MavMessageIdSerialControl:                      "SERIAL_CONTROL",
	MavMessageIdGpsRtk:                             "GPS_RTK",
	MavMessageIdGps2Rtk:                            "GPS2_RTK",
	MavMessageIdScaledImu3:                         "SCALED_IMU3",
	MavMessageIdDataTransmissionHandshake:          "DATA_TRANSMISSION_HANDSHAKE",
	MavMessageIdEncapsulatedData:                   "ENCAPSULATED_DATA",
	MavMessageIdDistanceSensor:                     "DISTANCE_SENSOR",
	MavMessageIdTerrainRequest:                     "TERRAIN_REQUEST",
	MavMessageIdTerrainData:                        "TERRAIN_DATA",
	MavMessageIdTerrainCheck:                       "TERRAIN_CHECK",
	MavMessageIdTerrainReport:                      "TERRAIN_REPORT",
	MavMessageIdScaledPressure2:                    "SCALED_PRESSURE2",
	MavMessageIdAttPosMocap:                        "ATT_POS_MOCAP",
	MavMessageIdSetActuatorControlTarget:           "SET_ACTUATOR_CONTROL_TARGET",
	MavMessageIdActuatorControlTarget:              "ACTUATOR_CONTROL_TARGET",
	MavMessageIdAltitude:                           "ALTITUDE",
	MavMessageIdResourceRequest:                    "RESOURCE_REQUEST",
	MavMessageIdScaledPressure3:                    "SCALED_PRESSURE3",
	MavMessageIdFollowTarget:                       "FOLLOW_TARGET",
	MavMessageIdControlSystemState:                 "CONTROL_SYSTEM_STATE",
	MavMessageIdBatteryStatus:                      "BATTERY_STATUS",
	MavMessageIdAutopilotVersion:                   "AUTOPILOT_VERSION",
	MavMessageIdLandingTarget:                      "LANDING_TARGET",
	MavMessageIdSensorOffsets:                      "SENSOR_OFFSETS",
	MavMessageIdSetMagOffsets:                      "SET_MAG_OFFSETS",
	MavMessageIdMeminfo:                            "MEMINFO",
	MavMessageIdApAdc:                              "AP_ADC",
	MavMessageIdDigicamConfigure:                   "DIGICAM_CONFIGURE",
	MavMessageIdDigicamControl:                     "DIGICAM_CONTROL",
	MavMessageIdMountConfigure:                     "MOUNT_CONFIGURE",
	MavMessageIdMountControl:                       "MOUNT_CONTROL",
	MavMessageIdMountStatus:                        "MOUNT_STATUS",
	MavMessageIdFencePoint:                         "FENCE_POINT",
	MavMessageIdFenceFetchPoint:                    "FENCE_FETCH_POINT",
	MavMessageIdFenceStatus:                        "FENCE_STATUS",
	MavMessageIdAhrs:                               "AHRS",
	MavMessageIdSimstate:                           "SIMSTATE",
	MavMessageIdHwstatus:                           "HWSTATUS",
	MavMessageIdRadio:                              "RADIO",
	MavMessageIdLimitsStatus:                       "LIMITS_STATUS",
	MavMessageIdWind:                               "WIND",
	MavMessageIdData16:                             "DATA16",
	MavMessageIdData32:                             "DATA32",
	MavMessageIdData64:                             "DATA64",
	MavMessageIdData96:                             "DATA96",
	MavMessageIdRangefinder:                        "RANGEFINDER",
	MavMessageIdAirspeedAutocal:                    "AIRSPEED_AUTOCAL",
	MavMessageIdRallyPoint:                         "RALLY_POINT",
	MavMessageIdRallyFetchPoint:                    "RALLY_FETCH_POINT",
	MavMessageIdCompassCalibrationProgress:         "COMPASS_CALIBRATION_PROGRESS",
	MavMessageIdMagCalReport:                       "MAG_CAL_REPORT",
	MavMessageIdEkfStatusReport:                    "EKF_STATUS_REPORT",
	MavMessageIdPidTuning:                          "PID_TUNING",
	MavMessageIdDeepstall:                          "DEEPSTALL",
	MavMessageIdGimbalReport:                       "GIMBAL_REPORT",
	MavMessageIdGimbalControl:                      "GIMBAL_CONTROL",
	MavMessageIdGimbalTorqueCmdReport:              "GIMBAL_TORQUE_CMD_REPORT",
	MavMessageIdGpsInput:                           "GPS_INPUT",
	MavMessageIdGpsRtcmData:                        "GPS_RTCM_DATA",
	MavMessageIdHighLatency:                        "HIGH_LATENCY",
	MavMessageIdHighLatency2:                       "HIGH_LATENCY2",
	MavMessageIdVibration:                          "VIBRATION",
	MavMessageIdHomePosition:                       "HOME_POSITION",
	MavMessageIdSetHomePosition:                    "SET_HOME_POSITION",
	MavMessageIdMessageInterval:                    "MESSAGE_INTERVAL",
	MavMessageIdExtendedSysState:                   "EXTENDED_SYS_STATE",
	MavMessageIdAdsbVehicle:                        "ADSB_VEHICLE",
	MavMessageIdCollision:                          "COLLISION",
	MavMessageIdV2Extension:                        "V2_EXTENSION",
	MavMessageIdMemoryVect:                         "MEMORY_VECT",
	MavMessageIdDebugVect:                          "DEBUG_VECT",
	MavMessageIdNamedValueFloat:                    "NAMED_VALUE_FLOAT",
	MavMessageIdNamedValueInt:                      "NAMED_VALUE_INT",
	MavMessageIdStatustext:                         "STATUSTEXT",
	MavMessageIdDebug:                              "DEBUG",
	MavMessageIdSetupSigning:                       "SETUP_SIGNING",
	MavMessageIdButtonChange:                       "BUTTON_CHANGE",
	MavMessageIdPlayTune:                           "PLAY_TUNE",
	MavMessageIdCameraInformation:                  "CAMERA_INFORMATION",
	MavMessageIdCameraSettings:                     "CAMERA_SETTINGS",
	MavMessageIdStorageInformation:                 "STORAGE_INFORMATION",
	MavMessageIdCameraCaptureStatus:                "CAMERA_CAPTURE_STATUS",
	MavMessageIdCameraImageCaptured:                "CAMERA_IMAGE_CAPTURED",
	MavMessageIdFlightInformation:                  "FLIGHT_INFORMATION",
	MavMessageIdMountOrientation:                   "MOUNT_ORIENTATION",
	MavMessageIdLoggingData:                        "LOGGING_DATA",
	MavMessageIdLoggingDataAcked:                   "LOGGING_DATA_ACKED",
	MavMessageIdLoggingAck:                         "LOGGING_ACK",
	MavMessageIdVideoStreamInformation:             "VIDEO_STREAM_INFORMATION",
	MavMessageIdVideoStreamStatus:                  "VIDEO_STREAM_STATUS",
	MavMessageIdCameraFovStatus:                    "CAMERA_FOV_STATUS",
	MavMessageIdCameraTrackingImageStatus:          "CAMERA_TRACKING_IMAGE_STATUS",
	MavMessageIdCameraTrackingGeoStatus:            "CAMERA_TRACKING_GEO_STATUS",
	MavMessageIdCameraThermalRange:                 "CAMERA_THERMAL_RANGE",
	MavMessageIdGimbalManagerInformation:           "GIMBAL_MANAGER_INFORMATION",
	MavMessageIdGimbalManagerStatus:                "GIMBAL_MANAGER_STATUS",
	MavMessageIdGimbalManagerSetAttitude:           "GIMBAL_MANAGER_SET_ATTITUDE",
	MavMessageIdGimbalManagerSetPitchyaw:           "GIMBAL_MANAGER_SET_PITCHYAW",
	MavMessageIdGimbalManagerSetManualControl:      "GIMBAL_MANAGER_SET_MANUAL_CONTROL",
	MavMessageIdGimbalDeviceInformation:            "GIMBAL_DEVICE_INFORMATION",
	MavMessageIdGimbalDeviceSetAttitude:            "GIMBAL_DEVICE_SET_ATTITUDE",
	MavMessageIdGimbalDeviceAttitudeStatus:         "GIMBAL_DEVICE_ATTITUDE_STATUS",
	MavMessageIdAutopilotStateForGimbalDevice:      "AUTOPILOT_STATE_FOR_GIMBAL_DEVICE",
	MavMessageIdEfiStatus:                          "EFI_STATUS",
	MavMessageIdEstimatorStatus:                    "ESTIMATOR_STATUS",
	MavMessageIdWindCov:                            "WIND_COV",
	MavMessageIdEscInfo:                            "ESC_INFO",
	MavMessageIdEscStatus:                          "ESC_STATUS",
	MavMessageIdAirspeed:                           "AIRSPEED",
	MavMessageIdWifiConfigAp:                       "WIFI_CONFIG_AP",
	MavMessageIdProtocolVersion:                    "PROTOCOL_VERSION",
	MavMessageIdAisVessel:                          "AIS_VESSEL",
	MavMessageIdUavcanNodeStatus:                   "UAVCAN_NODE_STATUS",
	MavMessageIdUavcanNodeInfo:                     "UAVCAN_NODE_INFO",
	MavMessageIdParamExtRequestRead:                "PARAM_EXT_REQUEST_READ",
	MavMessageIdParamExtRequestList:                "PARAM_EXT_REQUEST_LIST",
	MavMessageIdParamExtValue:                      "PARAM_EXT_VALUE",
	MavMessageIdParamExtSet:                        "PARAM_EXT_SET",
	MavMessageIdParamExtAck:                        "PARAM_EXT_ACK",
	MavMessageIdCellularStatus:                     "CELLULAR_STATUS",
	MavMessageIdIsbdLinkStatus:                     "ISBD_LINK_STATUS",
	MavMessageIdCellularConfig:                     "CELLULAR_CONFIG",
	MavMessageIdRawRpm:                             "RAW_RPM",
	MavMessageIdObstacleDistance:                   "OBSTACLE_DISTANCE",
	MavMessageIdOdometry:                           "ODOMETRY",
	MavMessageIdTrajectoryRepresentationWaypoints:  "TRAJECTORY_REPRESENTATION_WAYPOINTS",
	MavMessageIdTrajectoryRepresentationBezier:     "TRAJECTORY_REPRESENTATION_BEZIER",
	MavMessageIdUtmGlobalPosition:                  "UTM_GLOBAL_POSITION",
	MavMessageIdParamError:                         "PARAM_ERROR",
	MavMessageIdDebugFloatArray:                    "DEBUG_FLOAT_ARRAY",
	MavMessageIdOrbitExecutionStatus:               "ORBIT_EXECUTION_STATUS",
	MavMessageIdFigureEightExecutionStatus:         "FIGURE_EIGHT_EXECUTION_STATUS",
	MavMessageIdSmartBatteryInfo:                   "SMART_BATTERY_INFO",
	MavMessageIdFuelStatus:                         "FUEL_STATUS",
	MavMessageIdBatteryInfo:                        "BATTERY_INFO",
	MavMessageIdGeneratorStatus:                    "GENERATOR_STATUS",
	MavMessageIdActuatorOutputStatus:               "ACTUATOR_OUTPUT_STATUS",
	MavMessageIdTimeEstimateToTarget:               "TIME_ESTIMATE_TO_TARGET",
	MavMessageIdTunnel:                             "TUNNEL",
	MavMessageIdOnboardComputerStatus:              "ONBOARD_COMPUTER_STATUS",
	MavMessageIdComponentInformation:               "COMPONENT_INFORMATION",
	MavMessageIdComponentInformationBasic:          "COMPONENT_INFORMATION_BASIC",
	MavMessageIdComponentMetadata:                  "COMPONENT_METADATA",
	MavMessageIdComponentMetadataV2:                "COMPONENT_METADATA_V2",
	MavMessageIdPlayTuneV2:                         "PLAY_TUNE_V2",
	MavMessageIdSupportedTunes:                     "SUPPORTED_TUNES",
	MavMessageIdEvent:                              "EVENT",
	MavMessageIdCurrentEventSequence:               "CURRENT_EVENT_SEQUENCE",
	MavMessageIdRequestEvent:                       "REQUEST_EVENT",
	MavMessageIdResponseEventError:                 "RESPONSE_EVENT_ERROR",
	MavMessageIdCanFrame:                           "CAN_FRAME",
	MavMessageIdCanfdFrame:                         "CANFD_FRAME",
	MavMessageIdCanFilterModify:                    "CAN_FILTER_MODIFY",
	MavMessageIdWheelDistance:                      "WHEEL_DISTANCE",
	MavMessageIdWinchStatus:                        "WINCH_STATUS",
	MavMessageIdOpenDroneIdBasicId:                 "OPEN_DRONE_ID_BASIC_ID",
	MavMessageIdOpenDroneIdLocation:                "OPEN_DRONE_ID_LOCATION",
	MavMessageIdOpenDroneIdAuthentication:          "OPEN_DRONE_ID_AUTHENTICATION",
	MavMessageIdOpenDroneIdSelfId:                  "OPEN_DRONE_ID_SELF_ID",
	MavMessageIdOpenDroneIdSystem:                  "OPEN_DRONE_ID_SYSTEM",
	MavMessageIdOpenDroneIdOperatorId:              "OPEN_DRONE_ID_OPERATOR_ID",
	MavMessageIdOpenDroneIdMessagePack:             "OPEN_DRONE_ID_MESSAGE_PACK",
	MavMessageIdOpenDroneIdArmStatus:               "OPEN_DRONE_ID_ARM_STATUS",
	MavMessageIdOpenDroneIdSystemUpdate:            "OPEN_DRONE_ID_SYSTEM_UPDATE",
	MavMessageIdHygrometerSensor:                   "HYGROMETER_SENSOR",
	MavMessageIdCurrentMode:                        "CURRENT_MODE",
	MavMessageIdAvailableModes:                     "AVAILABLE_MODES",
	MavMessageIdAvailableModesMonitor:              "AVAILABLE_MODES_MONITOR",
	MavMessageIdIlluminatorStatus:                  "ILLUMINATOR_STATUS",
}

var stringToMavMessageId = map[string]MavMessageId{
	"HEARTBEAT":                               MavMessageIdHeartbeat,
	"SYS_STATUS":                              MavMessageIdSysStatus,
	"SYSTEM_TIME":                             MavMessageIdSystemTime,
	"PING":                                    MavMessageIdPing,
	"CHANGE_OPERATOR_CONTROL":                 MavMessageIdChangeOperatorControl,
	"CHANGE_OPERATOR_CONTROL_ACK":             MavMessageIdChangeOperatorControlAck,
	"AUTH_KEY":                                MavMessageIdAuthKey,
	"LINK_NODE_STATUS":                        MavMessageIdLinkNodeStatus,
	"SET_MODE":                                MavMessageIdSetMode,
	"PARAM_REQUEST_READ":                      MavMessageIdParamRequestRead,
	"PARAM_REQUEST_LIST":                      MavMessageIdParamRequestList,
	"PARAM_VALUE":                             MavMessageIdParamValue,
	"PARAM_SET":                               MavMessageIdParamSet,
	"GPS_RAW_INT":                             MavMessageIdGpsRawInt,
	"GPS_STATUS":                              MavMessageIdGpsStatus,
	"SCALED_IMU":                              MavMessageIdScaledImu,
	"RAW_IMU":                                 MavMessageIdRawImu,
	"RAW_PRESSURE":                            MavMessageIdRawPressure,
	"SCALED_PRESSURE":                         MavMessageIdScaledPressure,
	"ATTITUDE":                                MavMessageIdAttitude,
	"ATTITUDE_QUATERNION":                     MavMessageIdAttitudeQuaternion,
	"LOCAL_POSITION_NED":                      MavMessageIdLocalPositionNed,
	"GLOBAL_POSITION_INT":                     MavMessageIdGlobalPositionInt,
	"RC_CHANNELS_SCALED":                      MavMessageIdRcChannelsScaled,
	"RC_CHANNELS_RAW":                         MavMessageIdRcChannelsRaw,
	"SERVO_OUTPUT_RAW":                        MavMessageIdServoOutputRaw,
	"MISSION_REQUEST_PARTIAL_LIST":            MavMessageIdMissionRequestPartialList,
	"MISSION_WRITE_PARTIAL_LIST":              MavMessageIdMissionWritePartialList,
	"MISSION_ITEM":                            MavMessageIdMissionItem,
	"MISSION_REQUEST":                         MavMessageIdMissionRequest,
	"MISSION_SET_CURRENT":                     MavMessageIdMissionSetCurrent,
	"MISSION_CURRENT":                         MavMessageIdMissionCurrent,
	"MISSION_REQUEST_LIST":                    MavMessageIdMissionRequestList,
	"MISSION_COUNT":                           MavMessageIdMissionCount,
	"MISSION_CLEAR_ALL":                       MavMessageIdMissionClearAll,
	"MISSION_ITEM_REACHED":                    MavMessageIdMissionItemReached,
	"MISSION_ACK":                             MavMessageIdMissionAck,
	"SET_GPS_GLOBAL_ORIGIN":                   MavMessageIdSetGpsGlobalOrigin,
	"GPS_GLOBAL_ORIGIN":                       MavMessageIdGpsGlobalOrigin,
	"PARAM_MAP_RC":                            MavMessageIdParamMapRc,
	"MISSION_REQUEST_INT":                     MavMessageIdMissionRequestInt,
	"SAFETY_SET_ALLOWED_AREA":                 MavMessageIdSafetySetAllowedArea,
	"SAFETY_ALLOWED_AREA":                     MavMessageIdSafetyAllowedArea,
	"ATTITUDE_QUATERNION_COV":                 MavMessageIdAttitudeQuaternionCov,
	"NAV_CONTROLLER_OUTPUT":                   MavMessageIdNavControllerOutput,
	"GLOBAL_POSITION_INT_COV":                 MavMessageIdGlobalPositionIntCov,
	"LOCAL_POSITION_NED_COV":                  MavMessageIdLocalPositionNedCov,
	"RC_CHANNELS":                             MavMessageIdRcChannels,
	"REQUEST_DATA_STREAM":                     MavMessageIdRequestDataStream,
	"DATA_STREAM":                             MavMessageIdDataStream,
	"MANUAL_CONTROL":                          MavMessageIdManualControl,
	"RC_CHANNELS_OVERRIDE":                    MavMessageIdRcChannelsOverride,
	"MISSION_ITEM_INT":                        MavMessageIdMissionItemInt,
	"VFR_HUD":                                 MavMessageIdVfrHud,
	"COMMAND_INT":                             MavMessageIdCommandInt,
	"COMMAND_LONG":                            MavMessageIdCommandLong,
	"COMMAND_ACK":                             MavMessageIdCommandAck,
	"COMMAND_CANCEL":                          MavMessageIdCommandCancel,
	"MANUAL_SETPOINT":                         MavMessageIdManualSetpoint,
	"SET_ATTITUDE_TARGET":                     MavMessageIdSetAttitudeTarget,
	"ATTITUDE_TARGET":                         MavMessageIdAttitudeTarget,
	"SET_POSITION_TARGET_LOCAL_NED":           MavMessageIdSetPositionTargetLocalNed,
	"POSITION_TARGET_LOCAL_NED":               MavMessageIdPositionTargetLocalNed,
	"SET_POSITION_TARGET_GLOBAL_INT":          MavMessageIdSetPositionTargetGlobalInt,
	"POSITION_TARGET_GLOBAL_INT":              MavMessageIdPositionTargetGlobalInt,
	"POSITION_TARGET_GLOBAL_INT_REL_HOME":     MavMessageIdPositionTargetGlobalIntRelHome,
	"LOCAL_POSITION_NED_SYSTEM_GLOBAL_OFFSET": MavMessageIdLocalPositionNedSystemGlobalOffset,
	"HIL_STATE":                               MavMessageIdHilState,
	"HIL_CONTROLS":                            MavMessageIdHilControls,
	"HIL_RC_INPUTS_RAW":                       MavMessageIdHilRcInputsRaw,
	"HIL_ACTUATOR_CONTROLS":                   MavMessageIdHilActuatorControls,
	"OPTICAL_FLOW":                            MavMessageIdOpticalFlow,
	"GLOBAL_VISION_POSITION_ESTIMATE":         MavMessageIdGlobalVisionPositionEstimate,
	"VISION_POSITION_ESTIMATE":                MavMessageIdVisionPositionEstimate,
	"VISION_SPEED_ESTIMATE":                   MavMessageIdVisionSpeedEstimate,
	"VICON_POSITION_ESTIMATE":                 MavMessageIdViconPositionEstimate,
	"HIGHRES_IMU":                             MavMessageIdHighresImu,
	"OPTICAL_FLOW_RAD":                        MavMessageIdOpticalFlowRad,
	"HIL_SENSOR":                              MavMessageIdHilSensor,
	"SIM_STATE":                               MavMessageIdSimState,
	"RADIO_STATUS":                            MavMessageIdRadioStatus,
	"FILE_TRANSFER_PROTOCOL":                  MavMessageIdFileTransferProtocol,
	"TIMESYNC":                                MavMessageIdTimesync,
	"CAMERA_TRIGGER":                          MavMessageIdCameraTrigger,
	"HIL_GPS":                                 MavMessageIdHilGps,
	"HIL_OPTICAL_FLOW":                        MavMessageIdHilOpticalFlow,
	"HIL_STATE_QUATERNION":                    MavMessageIdHilStateQuaternion,
	"SCALED_IMU2":                             MavMessageIdScaledImu2,
	"LOG_REQUEST_LIST":                        MavMessageIdLogRequestList,
	"LOG_ENTRY":                               MavMessageIdLogEntry,
	"LOG_REQUEST_DATA":                        MavMessageIdLogRequestData,
	"LOG_DATA":                                MavMessageIdLogData,
	"LOG_ERASE":                               MavMessageIdLogErase,
	"LOG_REQUEST_END":                         MavMessageIdLogRequestEnd,
	"GPS_INJECT_DATA":                         MavMessageIdGpsInjectData,
	"GPS2_RAW":                                MavMessageIdGps2Raw,
	"POWER_STATUS":                            MavMessageIdPowerStatus,
	"SERIAL_CONTROL":                          MavMessageIdSerialControl,
	"GPS_RTK":                                 MavMessageIdGpsRtk,
	"GPS2_RTK":                                MavMessageIdGps2Rtk,
	"SCALED_IMU3":                             MavMessageIdScaledImu3,
	"DATA_TRANSMISSION_HANDSHAKE":             MavMessageIdDataTransmissionHandshake,
	"ENCAPSULATED_DATA":                       MavMessageIdEncapsulatedData,
	"DISTANCE_SENSOR":                         MavMessageIdDistanceSensor,
	"TERRAIN_REQUEST":                         MavMessageIdTerrainRequest,
	"TERRAIN_DATA":                            MavMessageIdTerrainData,
	"TERRAIN_CHECK":                           MavMessageIdTerrainCheck,
	"TERRAIN_REPORT":                          MavMessageIdTerrainReport,
	"SCALED_PRESSURE2":                        MavMessageIdScaledPressure2,
	"ATT_POS_MOCAP":                           MavMessageIdAttPosMocap,
	"SET_ACTUATOR_CONTROL_TARGET":             MavMessageIdSetActuatorControlTarget,
	"ACTUATOR_CONTROL_TARGET":                 MavMessageIdActuatorControlTarget,
	"ALTITUDE":                                MavMessageIdAltitude,
	"RESOURCE_REQUEST":                        MavMessageIdResourceRequest,
	"SCALED_PRESSURE3":                        MavMessageIdScaledPressure3,
	"FOLLOW_TARGET":                           MavMessageIdFollowTarget,
	"CONTROL_SYSTEM_STATE":                    MavMessageIdControlSystemState,
	"BATTERY_STATUS":                          MavMessageIdBatteryStatus,
	"AUTOPILOT_VERSION":                       MavMessageIdAutopilotVersion,
	"LANDING_TARGET":                          MavMessageIdLandingTarget,
	"SENSOR_OFFSETS":                          MavMessageIdSensorOffsets,
	"SET_MAG_OFFSETS":                         MavMessageIdSetMagOffsets,
	"MEMINFO":                                 MavMessageIdMeminfo,
	"AP_ADC":                                  MavMessageIdApAdc,
	"DIGICAM_CONFIGURE":                       MavMessageIdDigicamConfigure,
	"DIGICAM_CONTROL":                         MavMessageIdDigicamControl,
	"MOUNT_CONFIGURE":                         MavMessageIdMountConfigure,
	"MOUNT_CONTROL":                           MavMessageIdMountControl,
	"MOUNT_STATUS":                            MavMessageIdMountStatus,
	"FENCE_POINT":                             MavMessageIdFencePoint,
	"FENCE_FETCH_POINT":                       MavMessageIdFenceFetchPoint,
	"FENCE_STATUS":                            MavMessageIdFenceStatus,
	"AHRS":                                    MavMessageIdAhrs,
	"SIMSTATE":                                MavMessageIdSimstate,
	"HWSTATUS":                                MavMessageIdHwstatus,
	"RADIO":                                   MavMessageIdRadio,
	"LIMITS_STATUS":                           MavMessageIdLimitsStatus,
	"WIND":                                    MavMessageIdWind,
	"DATA16":                                  MavMessageIdData16,
	"DATA32":                                  MavMessageIdData32,
	"DATA64":                                  MavMessageIdData64,
	"DATA96":                                  MavMessageIdData96,
	"RANGEFINDER":                             MavMessageIdRangefinder,
	"AIRSPEED_AUTOCAL":                        MavMessageIdAirspeedAutocal,
	"RALLY_POINT":                             MavMessageIdRallyPoint,
	"RALLY_FETCH_POINT":                       MavMessageIdRallyFetchPoint,
	"COMPASS_CALIBRATION_PROGRESS":            MavMessageIdCompassCalibrationProgress,
	"MAG_CAL_REPORT":                          MavMessageIdMagCalReport,
	"EKF_STATUS_REPORT":                       MavMessageIdEkfStatusReport,
	"PID_TUNING":                              MavMessageIdPidTuning,
	"DEEPSTALL":                               MavMessageIdDeepstall,
	"GIMBAL_REPORT":                           MavMessageIdGimbalReport,
	"GIMBAL_CONTROL":                          MavMessageIdGimbalControl,
	"GIMBAL_TORQUE_CMD_REPORT":                MavMessageIdGimbalTorqueCmdReport,
	"GPS_INPUT":                               MavMessageIdGpsInput,
	"GPS_RTCM_DATA":                           MavMessageIdGpsRtcmData,
	"HIGH_LATENCY":                            MavMessageIdHighLatency,
	"HIGH_LATENCY2":                           MavMessageIdHighLatency2,
	"VIBRATION":                               MavMessageIdVibration,
	"HOME_POSITION":                           MavMessageIdHomePosition,
	"SET_HOME_POSITION":                       MavMessageIdSetHomePosition,
	"MESSAGE_INTERVAL":                        MavMessageIdMessageInterval,
	"EXTENDED_SYS_STATE":                      MavMessageIdExtendedSysState,
	"ADSB_VEHICLE":                            MavMessageIdAdsbVehicle,
	"COLLISION":                               MavMessageIdCollision,
	"V2_EXTENSION":                            MavMessageIdV2Extension,
	"MEMORY_VECT":                             MavMessageIdMemoryVect,
	"DEBUG_VECT":                              MavMessageIdDebugVect,
	"NAMED_VALUE_FLOAT":                       MavMessageIdNamedValueFloat,
	"NAMED_VALUE_INT":                         MavMessageIdNamedValueInt,
	"STATUSTEXT":                              MavMessageIdStatustext,
	"DEBUG":                                   MavMessageIdDebug,
	"SETUP_SIGNING":                           MavMessageIdSetupSigning,
	"BUTTON_CHANGE":                           MavMessageIdButtonChange,
	"PLAY_TUNE":                               MavMessageIdPlayTune,
	"CAMERA_INFORMATION":                      MavMessageIdCameraInformation,
	"CAMERA_SETTINGS":                         MavMessageIdCameraSettings,
	"STORAGE_INFORMATION":                     MavMessageIdStorageInformation,
	"CAMERA_CAPTURE_STATUS":                   MavMessageIdCameraCaptureStatus,
	"CAMERA_IMAGE_CAPTURED":                   MavMessageIdCameraImageCaptured,
	"FLIGHT_INFORMATION":                      MavMessageIdFlightInformation,
	"MOUNT_ORIENTATION":                       MavMessageIdMountOrientation,
	"LOGGING_DATA":                            MavMessageIdLoggingData,
	"LOGGING_DATA_ACKED":                      MavMessageIdLoggingDataAcked,
	"LOGGING_ACK":                             MavMessageIdLoggingAck,
	"VIDEO_STREAM_INFORMATION":                MavMessageIdVideoStreamInformation,
	"VIDEO_STREAM_STATUS":                     MavMessageIdVideoStreamStatus,
	"CAMERA_FOV_STATUS":                       MavMessageIdCameraFovStatus,
	"CAMERA_TRACKING_IMAGE_STATUS":            MavMessageIdCameraTrackingImageStatus,
	"CAMERA_TRACKING_GEO_STATUS":              MavMessageIdCameraTrackingGeoStatus,
	"CAMERA_THERMAL_RANGE":                    MavMessageIdCameraThermalRange,
	"GIMBAL_MANAGER_INFORMATION":              MavMessageIdGimbalManagerInformation,
	"GIMBAL_MANAGER_STATUS":                   MavMessageIdGimbalManagerStatus,
	"GIMBAL_MANAGER_SET_ATTITUDE":             MavMessageIdGimbalManagerSetAttitude,
	"GIMBAL_MANAGER_SET_PITCHYAW":             MavMessageIdGimbalManagerSetPitchyaw,
	"GIMBAL_MANAGER_SET_MANUAL_CONTROL":       MavMessageIdGimbalManagerSetManualControl,
	"GIMBAL_DEVICE_INFORMATION":               MavMessageIdGimbalDeviceInformation,
	"GIMBAL_DEVICE_SET_ATTITUDE":              MavMessageIdGimbalDeviceSetAttitude,
	"GIMBAL_DEVICE_ATTITUDE_STATUS":           MavMessageIdGimbalDeviceAttitudeStatus,
	"AUTOPILOT_STATE_FOR_GIMBAL_DEVICE":       MavMessageIdAutopilotStateForGimbalDevice,
	"EFI_STATUS":                              MavMessageIdEfiStatus,
	"ESTIMATOR_STATUS":                        MavMessageIdEstimatorStatus,
	"WIND_COV":                                MavMessageIdWindCov,
	"ESC_INFO":                                MavMessageIdEscInfo,
	"ESC_STATUS":                              MavMessageIdEscStatus,
	"AIRSPEED":                                MavMessageIdAirspeed,
	"WIFI_CONFIG_AP":                          MavMessageIdWifiConfigAp,
	"PROTOCOL_VERSION":                        MavMessageIdProtocolVersion,
	"AIS_VESSEL":                              MavMessageIdAisVessel,
	"UAVCAN_NODE_STATUS":                      MavMessageIdUavcanNodeStatus,
	"UAVCAN_NODE_INFO":                        MavMessageIdUavcanNodeInfo,
	"PARAM_EXT_REQUEST_READ":                  MavMessageIdParamExtRequestRead,
	"PARAM_EXT_REQUEST_LIST":                  MavMessageIdParamExtRequestList,
	"PARAM_EXT_VALUE":                         MavMessageIdParamExtValue,
	"PARAM_EXT_SET":                           MavMessageIdParamExtSet,
	"PARAM_EXT_ACK":                           MavMessageIdParamExtAck,
	"CELLULAR_STATUS":                         MavMessageIdCellularStatus,
	"ISBD_LINK_STATUS":                        MavMessageIdIsbdLinkStatus,
	"CELLULAR_CONFIG":                         MavMessageIdCellularConfig,
	"RAW_RPM":                                 MavMessageIdRawRpm,
	"OBSTACLE_DISTANCE":                       MavMessageIdObstacleDistance,
	"ODOMETRY":                                MavMessageIdOdometry,
	"TRAJECTORY_REPRESENTATION_WAYPOINTS":     MavMessageIdTrajectoryRepresentationWaypoints,
	"TRAJECTORY_REPRESENTATION_BEZIER":        MavMessageIdTrajectoryRepresentationBezier,
	"UTM_GLOBAL_POSITION":                     MavMessageIdUtmGlobalPosition,
	"PARAM_ERROR":                             MavMessageIdParamError,
	"DEBUG_FLOAT_ARRAY":                       MavMessageIdDebugFloatArray,
	"ORBIT_EXECUTION_STATUS":                  MavMessageIdOrbitExecutionStatus,
	"FIGURE_EIGHT_EXECUTION_STATUS":           MavMessageIdFigureEightExecutionStatus,
	"SMART_BATTERY_INFO":                      MavMessageIdSmartBatteryInfo,
	"FUEL_STATUS":                             MavMessageIdFuelStatus,
	"BATTERY_INFO":                            MavMessageIdBatteryInfo,
	"GENERATOR_STATUS":                        MavMessageIdGeneratorStatus,
	"ACTUATOR_OUTPUT_STATUS":                  MavMessageIdActuatorOutputStatus,
	"TIME_ESTIMATE_TO_TARGET":                 MavMessageIdTimeEstimateToTarget,
	"TUNNEL":                                  MavMessageIdTunnel,
	"ONBOARD_COMPUTER_STATUS":                 MavMessageIdOnboardComputerStatus,
	"COMPONENT_INFORMATION":                   MavMessageIdComponentInformation,
	"COMPONENT_INFORMATION_BASIC":             MavMessageIdComponentInformationBasic,
	"COMPONENT_METADATA":                      MavMessageIdComponentMetadata,
	"COMPONENT_METADATA_V2":                   MavMessageIdComponentMetadataV2,
	"PLAY_TUNE_V2":                            MavMessageIdPlayTuneV2,
	"SUPPORTED_TUNES":                         MavMessageIdSupportedTunes,
	"EVENT":                                   MavMessageIdEvent,
	"CURRENT_EVENT_SEQUENCE":                  MavMessageIdCurrentEventSequence,
	"REQUEST_EVENT":                           MavMessageIdRequestEvent,
	"RESPONSE_EVENT_ERROR":                    MavMessageIdResponseEventError,
	"CAN_FRAME":                               MavMessageIdCanFrame,
	"CANFD_FRAME":                             MavMessageIdCanfdFrame,
	"CAN_FILTER_MODIFY":                       MavMessageIdCanFilterModify,
	"WHEEL_DISTANCE":                          MavMessageIdWheelDistance,
	"WINCH_STATUS":                            MavMessageIdWinchStatus,
	"OPEN_DRONE_ID_BASIC_ID":                  MavMessageIdOpenDroneIdBasicId,
	"OPEN_DRONE_ID_LOCATION":                  MavMessageIdOpenDroneIdLocation,
	"OPEN_DRONE_ID_AUTHENTICATION":            MavMessageIdOpenDroneIdAuthentication,
	"OPEN_DRONE_ID_SELF_ID":                   MavMessageIdOpenDroneIdSelfId,
	"OPEN_DRONE_ID_SYSTEM":                    MavMessageIdOpenDroneIdSystem,
	"OPEN_DRONE_ID_OPERATOR_ID":               MavMessageIdOpenDroneIdOperatorId,
	"OPEN_DRONE_ID_MESSAGE_PACK":              MavMessageIdOpenDroneIdMessagePack,
	"OPEN_DRONE_ID_ARM_STATUS":                MavMessageIdOpenDroneIdArmStatus,
	"OPEN_DRONE_ID_SYSTEM_UPDATE":             MavMessageIdOpenDroneIdSystemUpdate,
	"HYGROMETER_SENSOR":                       MavMessageIdHygrometerSensor,
	"CURRENT_MODE":                            MavMessageIdCurrentMode,
	"AVAILABLE_MODES":                         MavMessageIdAvailableModes,
	"AVAILABLE_MODES_MONITOR":                 MavMessageIdAvailableModesMonitor,
	"ILLUMINATOR_STATUS":                      MavMessageIdIlluminatorStatus,
}

// String
// Returns the string representation of the MavMessageId.
func (id MavMessageId) String() string {
	if str, ok := mavMessageIdStrings[id]; ok {
		return str
	}
	return fmt.Sprintf("UNKNOWN_MESSAGE_ID_%d", id)
}

// ParseMavMessageId
// Parses a string into a MavMessageId.
// Returns an error if the string does not correspond to a valid message ID.
func ParseMavMessageId(s string) (MavMessageId, error) {
	if id, ok := stringToMavMessageId[s]; ok {
		return id, nil
	}
	return 0, fmt.Errorf("invalid message ID: %s", s)
}
